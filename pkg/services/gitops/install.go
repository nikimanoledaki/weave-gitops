package gitops

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"github.com/weaveworks/weave-gitops/cmd/gitops/version"
	"github.com/weaveworks/weave-gitops/manifests"
	"github.com/weaveworks/weave-gitops/pkg/apputils"
	"github.com/weaveworks/weave-gitops/pkg/git"
	"github.com/weaveworks/weave-gitops/pkg/gitproviders"
	"github.com/weaveworks/weave-gitops/pkg/kube"
	"github.com/weaveworks/weave-gitops/pkg/services/automation"
	"github.com/weaveworks/weave-gitops/pkg/services/gitrepo"
)

type InstallParams struct {
	Namespace    string
	DryRun       bool
	AppConfigURL string
}

func (g *Gitops) Install(params InstallParams) ([]byte, error) {
	ctx := context.Background()
	status := g.kube.GetClusterStatus(ctx)

	switch status {
	case kube.FluxInstalled:
		return nil, errors.New("Weave GitOps does not yet support installation onto a cluster that is using Flux.\nPlease uninstall flux before proceeding:\n  $ flux uninstall")
	case kube.Unknown:
		return nil, errors.New("Weave GitOps cannot talk to the cluster")
	}

	// TODO apply these manifests instead of generating them again
	var fluxManifests []byte

	var err error

	if params.AppConfigURL != "" || params.DryRun {
		// We need to get the manifests to persist in the repo and
		// non-dry run install doesn't return them
		fluxManifests, err = g.flux.Install(params.Namespace, true)
		if err != nil {
			return fluxManifests, fmt.Errorf("error on flux install %w", err)
		}
	}

	if !params.DryRun {
		_, err = g.flux.Install(params.Namespace, false)
		if err != nil {
			return fluxManifests, fmt.Errorf("error on flux install %w", err)
		}
	}

	systemManifests := make(map[string][]byte, 3)
	systemManifests["gitops-runtime.yaml"] = fluxManifests
	systemManifests["wego-system.yaml"] = manifests.AppCRD

	if params.DryRun {
		fluxManifests = append(fluxManifests, manifests.AppCRD...)
	} else {

		for _, manifest := range manifests.Manifests {
			if err := g.kube.Apply(ctx, manifest, params.Namespace); err != nil {
				return nil, fmt.Errorf("could not apply manifest: %w", err)
			}
		}

		version := version.Version
		if os.Getenv("IS_TEST_ENV") != "" {
			version = "latest"
		}

		wegoAppDeploymentManifest, err := manifests.GenerateWegoAppDeploymentManifest(version)
		if err != nil {
			return nil, fmt.Errorf("error generating wego-app deployment, %w", err)
		}

		systemManifests["wego-app.yaml"] = wegoAppDeploymentManifest
		if err := g.kube.Apply(ctx, wegoAppDeploymentManifest, params.Namespace); err != nil {
			return nil, fmt.Errorf("could not apply wego-app deployment manifest: %w", err)
		}

		if params.AppConfigURL != "" {
			cname, err := g.kube.GetClusterName(ctx)
			if err != nil {
				g.logger.Warningf("Cluster name not found, using default : %v", err)

				cname = "default"
			}

			goatManifests, err := g.storeManifests(params, systemManifests, cname)
			if err != nil {
				return nil, fmt.Errorf("failed to store cluster manifests: %w", err)
			}

			g.logger.Actionf("Applying manifests to the cluster")
			// only apply the system manifests as the others will get picked up once flux is running
			if err := g.applyManifestsToK8s(ctx, params.Namespace, goatManifests); err != nil {
				return nil, fmt.Errorf("failed applying system manifests to cluster %s :%w", cname, err)
			}

		}
	}
	// TODO Existing install doesn't expect the systemManifests to be included in the list of manifests returned.
	// This will need to be changed.
	return fluxManifests, nil
}

func (g *Gitops) storeManifests(params InstallParams, systemManifests map[string][]byte, cname string) (map[string][]byte, error) {
	ctx := context.Background()

	normalizedURL, err := gitproviders.NewRepoURL(params.AppConfigURL)
	if err != nil {
		return nil, fmt.Errorf("failed to convert app config repo %q : %w", params.AppConfigURL, err)
	}

	configBranch, err := g.gitProvider.GetDefaultBranch(ctx, normalizedURL)
	if err != nil {
		return nil, fmt.Errorf("could not determine default branch for config repository: %q %w", params.AppConfigURL, err)
	}

	if g.gitClient == nil {
		authsvc, err := apputils.GetAuthService(ctx, normalizedURL, params.DryRun)
		if err != nil {
			return nil, fmt.Errorf("failed to create auth service for repo %q : %w", params.AppConfigURL, err)
		}

		g.gitClient, err = authsvc.CreateGitClient(ctx, normalizedURL, cname, params.Namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to create git client for repo %q : %w", params.AppConfigURL, err)
		}
	}

	//	remover, _, err := app.CloneRepo(g.gitClient, params.AppConfigURL, configBranch, params.DryRun)
	remover, _, err := gitrepo.CloneRepo(ctx, g.gitClient, normalizedURL, configBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to clone configuration repo: %w", err)
	}

	defer remover()

	manifests := make(map[string][]byte, 3)
	clusterPath := filepath.Join(git.WegoRoot, git.WegoClusterDir, cname)

	gitsource, sourceName, err := g.genSource(cname, configBranch, params.Namespace, normalizedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create source manifest: %w", err)
	}
	// filepath.Join doesn't support a starting "."
	prefixForFlux := func(s string) string {
		return "." + s
	}
	manifests["flux-source-resource.yaml"] = gitsource

	system, err := g.genKustomize(automation.ConstrainResourceName(fmt.Sprintf("%s-system", cname)), sourceName, configBranch,
		prefixForFlux(filepath.Join(".", clusterPath, git.WegoClusterOSWorkloadDir)), params)
	if err != nil {
		return nil, fmt.Errorf("failed to create system kustomization manifest: %w", err)
	}

	manifests["flux-system-kustomization-resource.yaml"] = system

	user, err := g.genKustomize(automation.ConstrainResourceName(fmt.Sprintf("%s-user", cname)), sourceName, configBranch,
		prefixForFlux(filepath.Join(".", clusterPath, git.WegoClusterUserWorloadDir)), params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user kustomization manifest: %w", err)
	}

	manifests["flux-user-kustomization-resource.yaml"] = user

	//TODO add handling for PRs
	// if !params.AutoMerge {
	//  if err := a.createPullRequestToRepo(info, info.Spec.ConfigURL, appHash, appSpec, appGoat, appSource); err != nil {
	//      return err
	//  }
	// } else {
	g.logger.Actionf("Writing manifests to disk")

	if err := g.writeManifestsToGit(filepath.Join(clusterPath, "system"), manifests); err != nil {
		return nil, fmt.Errorf("failed to write manifests: %w", err)
	}

	if err := g.writeManifestsToGit(filepath.Join(clusterPath, "system"), systemManifests); err != nil {
		return nil, fmt.Errorf("failed to write system manifests: %w", err)
	}
	// store a .keep file in the user dir
	userKeep := map[string][]byte{
		".keep": strconv.AppendQuote(nil, "# keep"),
	}
	if err := g.writeManifestsToGit(filepath.Join(clusterPath, "user"), userKeep); err != nil {
		return nil, fmt.Errorf("failed to write user manifests: %w", err)
	}

	//	return manifests, app.CommitAndPush(g.gitClient, "Add GitOps runtime manifests", params.DryRun, g.logger)
	return manifests, gitrepo.CommitAndPush(ctx, g.gitClient, "Add GitOps runtime manifests", g.logger)
}

func (g *Gitops) genSource(cname, branch string, namespace string, normalizedUrl gitproviders.RepoURL) ([]byte, string, error) {
	secretRef := automation.CreateRepoSecretName(cname, normalizedUrl).String()

	sourceManifest, err := g.flux.CreateSourceGit(secretRef, normalizedUrl, branch, secretRef, namespace)
	if err != nil {
		return nil, secretRef, fmt.Errorf("could not create git source for repo %s : %w", normalizedUrl.String(), err)
	}

	return sourceManifest, secretRef, nil
}

func (g *Gitops) genKustomize(name, cname, branch, path string, params InstallParams) ([]byte, error) {
	sourceManifest, err := g.flux.CreateKustomization(name, cname, path, params.Namespace)
	if err != nil {
		return nil, fmt.Errorf("could not create flux kustomization for path %q : %w", path, err)
	}

	return sourceManifest, nil
}

func (g *Gitops) writeManifestsToGit(path string, manifests map[string][]byte) error {
	for k, m := range manifests {
		if err := g.gitClient.Write(filepath.Join(path, k), m); err != nil {
			g.logger.Warningf("failed to write manifest %s : %v", k, err)
			return err
		}
	}

	return nil
}

func (g *Gitops) applyManifestsToK8s(ctx context.Context, namespace string, manifests map[string][]byte) error {
	for k, manifest := range manifests {
		if err := g.kube.Apply(ctx, manifest, namespace); err != nil {
			return fmt.Errorf("could not apply manifest %q : %w", k, err)
		}
	}

	return nil
}
