package clustersmngr

import (
	"context"
	"fmt"

	"github.com/weaveworks/weave-gitops/pkg/server/auth"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/cli-utils/pkg/flowcontrol"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type key int

const (
	// Clusters Client context key
	ClustersClientCtxKey key = iota
	// DefaultCluster name
	DefaultCluster = "Default"
	// ClientQPS is the QPS to use while creating the k8s clients
	ClientQPS = 1000
	// ClientBurst is the burst to use while creating the k8s clients
	ClientBurst = 2000
)

// Cluster defines a leaf cluster
type Cluster struct {
	// Name defines the cluster name
	Name string `yaml:"name"`
	// Server defines cluster api address
	Server string `yaml:"server"`

	// SecretRef defines secret name that holds the cluster Bearer Token
	SecretRef string `yaml:"secretRef"`
	// BearerToken cluster access token read from SecretRef
	BearerToken string

	// TLSConfig holds configuration for TLS connection with the cluster values read from SecretRef
	TLSConfig rest.TLSClientConfig
}

// ClusterNotFoundError cluster client can be found in the pool
type ClusterNotFoundError struct {
	Cluster string
}

func (e ClusterNotFoundError) Error() string {
	return fmt.Sprintf("cluster=%s not found", e.Cluster)
}

//ClusterFetcher fetches all leaf clusters
//counterfeiter:generate . ClusterFetcher
type ClusterFetcher interface {
	Fetch(ctx context.Context) ([]Cluster, error)
}

// ClientsPool stores all clients to the leaf clusters
//counterfeiter:generate . ClientsPool
type ClientsPool interface {
	Add(cfg ClusterClientConfig, cluster Cluster) error
	Clients() map[string]client.Client
	Client(cluster string) (client.Client, error)
}

type clientsPool struct {
	clients map[string]client.Client
	scheme  *apiruntime.Scheme
}

type ClusterClientConfig func(Cluster) *rest.Config

func ClientConfigWithUser(user *auth.UserPrincipal) ClusterClientConfig {
	return func(cluster Cluster) *rest.Config {
		config := &rest.Config{
			Host:            cluster.Server,
			BearerToken:     cluster.BearerToken,
			TLSClientConfig: cluster.TLSConfig,
			Impersonate: rest.ImpersonationConfig{
				UserName: user.ID,
				Groups:   user.Groups,
			},
		}

		enabled, err := flowcontrol.IsEnabled(context.Background(), config)
		if err == nil && enabled {
			// Enabled & negative QPS and Burst indicate that the client would use the rate limit set by the server.
			// Ref: https://github.com/kubernetes/kubernetes/blob/v1.24.0/staging/src/k8s.io/client-go/rest/config.go#L354-L364
			config.QPS = -1
			config.Burst = -1

			return config
		}

		config.QPS = ClientQPS
		config.Burst = ClientBurst

		return config
	}
}

// NewClustersClientsPool initializes a new ClientsPool
func NewClustersClientsPool(scheme *apiruntime.Scheme) ClientsPool {
	return &clientsPool{
		clients: map[string]client.Client{},
		scheme:  scheme,
	}
}

// Add adds a cluster client to the clients pool with the given user impersonation
func (cp *clientsPool) Add(cfg ClusterClientConfig, cluster Cluster) error {
	config := cfg(cluster)

	leafClient, err := client.New(config, client.Options{
		Scheme: cp.scheme,
	})
	if err != nil {
		return fmt.Errorf("failed to create leaf client: %w", err)
	}

	cp.clients[cluster.Name] = leafClient

	return nil
}

// Clients returns the clusters clients
func (cp *clientsPool) Clients() map[string]client.Client {
	return cp.clients
}

// Client returns the client for the given cluster
func (cp *clientsPool) Client(name string) (client.Client, error) {
	if c, found := cp.clients[name]; found && c != nil {
		return c, nil
	}

	return nil, ClusterNotFoundError{Cluster: name}
}
