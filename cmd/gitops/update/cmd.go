package update

import (
	"github.com/go-resty/resty/v2"
	"github.com/weaveworks/weave-gitops/cmd/gitops/update/profiles"

	"github.com/spf13/cobra"
)

func UpdateCommand(endpoint, username, password *string, client *resty.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a Weave GitOps resource",
		Example: `
	# Update a profile that is installed on a cluster
	gitops update profile --name=podinfo --cluster=prod --config-repo=ssh://git@github.com/owner/config-repo.git  --version=1.0.0
		`,
	}

	cmd.AddCommand(profiles.UpdateCommand(endpoint, username, password, client))

	return cmd
}
