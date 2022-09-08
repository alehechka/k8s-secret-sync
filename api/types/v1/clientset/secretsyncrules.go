package clientset

import (
	"context"
	"time"

	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const resource = "secretsyncrules"

// SecretSyncRuleGetter has a method to return a SecretSyncRuleInterface.
type SecretSyncRuleGetter interface {
	SecretSyncRules() SecretSyncRuleInterface
}

// secretSyncRules implements SecretSyncRuleInterface
type secretSyncRules struct {
	client rest.Interface
}

// newSecretSyncRules returns a SecretSyncRules
func newSecretSyncRules(c *KubeSecretSyncClientset) *secretSyncRules {
	return &secretSyncRules{
		client: c.client,
	}
}

// SecretSyncRuleInterface has methods to work with SecretSyncRule resources.
type SecretSyncRuleInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*typesv1.SecretSyncRuleList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*typesv1.SecretSyncRule, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

func (c *secretSyncRules) List(ctx context.Context, opts metav1.ListOptions) (*typesv1.SecretSyncRuleList, error) {
	result := typesv1.SecretSyncRuleList{}
	err := c.client.
		Get().
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *secretSyncRules) Get(ctx context.Context, name string, opts metav1.GetOptions) (*typesv1.SecretSyncRule, error) {
	result := typesv1.SecretSyncRule{}
	err := c.client.
		Get().
		Resource(resource).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *secretSyncRules) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.
		Get().
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}
