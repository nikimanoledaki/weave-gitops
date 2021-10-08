package gitproviders

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("detectGitProviderFromUrl", func(input string, expected GitProviderName) {
	result, err := detectGitProviderFromUrl(input)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).To(Equal(expected))
},
	Entry("ssh+github", "ssh://git@github.com/weaveworks/weave-gitops.git", GitProviderGitHub),
	Entry("ssh+gitlab", "ssh://git@gitlab.com/weaveworks/weave-gitops.git", GitProviderGitLab),
)

var _ = Describe("get owner from url", func() {
	DescribeTable("getOwnerFromUrl", func(normalizedUrl string, providerName GitProviderName, expected string) {
		u, err := url.Parse(normalizedUrl)
		Expect(err).NotTo(HaveOccurred())
		result, err := getOwnerFromUrl(*u, providerName)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(expected))
	},
		Entry("github", "ssh://git@github.com/weaveworks/weave-gitops.git", GitProviderGitHub, "weaveworks"),
		Entry("gitlab", "ssh://git@gitlab.com/weaveworks/weave-gitops.git", GitProviderGitLab, "weaveworks"),
		Entry("gitlab with subgroup", "ssh://git@gitlab.com/weaveworks/sub_group/weave-gitops.git", GitProviderGitLab, "weaveworks/sub_group"),
	)

	It("missing owner", func() {
		normalizedUrl := "ssh://git@gitlab.com/weave-gitops.git"
		u, err := url.Parse(normalizedUrl)
		Expect(err).NotTo(HaveOccurred())
		_, err = getOwnerFromUrl(*u, GitProviderGitLab)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("could not get owner from url ssh://git@gitlab.com/weave-gitops.git"))
	})

	It("empty url", func() {
		normalizedUrl := ""
		u, err := url.Parse(normalizedUrl)
		Expect(err).NotTo(HaveOccurred())
		_, err = getOwnerFromUrl(*u, GitProviderGitLab)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("could not get owner from url "))
	})

	It("subgroup in a subgroup", func() {
		normalizedUrl := "ssh://git@gitlab.com/weaveworks/sub_group/another_sub_group/weave-gitops.git"
		u, err := url.Parse(normalizedUrl)
		Expect(err).NotTo(HaveOccurred())
		_, err = getOwnerFromUrl(*u, GitProviderGitLab)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("a subgroup in a subgroup is not currently supported"))
	})
})

type expectedRepoURL struct {
	s        string
	owner    string
	name     string
	provider GitProviderName
	protocol RepositoryURLProtocol
}

var _ = DescribeTable("NormalizedRepoURL", func(input string, expected expectedRepoURL) {
	result, err := NewNormalizedRepoURL(input)
	Expect(err).NotTo(HaveOccurred())

	Expect(result.String()).To(Equal(expected.s))
	u, err := url.Parse(expected.s)
	Expect(err).NotTo(HaveOccurred())
	Expect(result.URL()).To(Equal(u))
	Expect(result.Owner()).To(Equal(expected.owner))
	Expect(result.Provider()).To(Equal(expected.provider))
	Expect(result.Protocol()).To(Equal(expected.protocol))
},
	Entry("github git clone style", "git@github.com:someuser/podinfo.git", expectedRepoURL{
		s:        "ssh://git@github.com/someuser/podinfo.git",
		owner:    "someuser",
		name:     "podinfo",
		provider: GitProviderGitHub,
		protocol: RepositoryURLProtocolSSH,
	}),
	Entry("github url style", "ssh://git@github.com/someuser/podinfo.git", expectedRepoURL{
		s:        "ssh://git@github.com/someuser/podinfo.git",
		owner:    "someuser",
		name:     "podinfo",
		provider: GitProviderGitHub,
		protocol: RepositoryURLProtocolSSH,
	}),
	Entry("github https", "https://github.com/someuser/podinfo.git", expectedRepoURL{
		s:        "ssh://git@github.com/someuser/podinfo.git",
		owner:    "someuser",
		name:     "podinfo",
		provider: GitProviderGitHub,
		protocol: RepositoryURLProtocolSSH,
	}),
	Entry("gitlab git clone style", "git@gitlab.com:someuser/podinfo.git", expectedRepoURL{
		s:        "ssh://git@gitlab.com/someuser/podinfo.git",
		owner:    "someuser",
		name:     "podinfo",
		provider: GitProviderGitLab,
		protocol: RepositoryURLProtocolSSH,
	}),
	Entry("gitlab https", "https://gitlab.com/someuser/podinfo.git", expectedRepoURL{
		s:        "ssh://git@gitlab.com/someuser/podinfo.git",
		owner:    "someuser",
		name:     "podinfo",
		provider: GitProviderGitLab,
		protocol: RepositoryURLProtocolSSH,
	}),
)