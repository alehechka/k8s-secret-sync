package clientset

import (
	"context"
	"time"

	v1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const resource = "secretsyncrules"

type SecretSyncRuleInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1.SecretSyncRuleList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1.SecretSyncRule, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type secretSyncRuleClient struct {
	restClient rest.Interface
	ns         string
}

func (c *secretSyncRuleClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.SecretSyncRuleList, error) {
	result := v1.SecretSyncRuleList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *secretSyncRuleClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.SecretSyncRule, error) {
	result := v1.SecretSyncRule{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resource).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *secretSyncRuleClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}
