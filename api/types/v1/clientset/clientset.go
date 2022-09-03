package clientset

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type KubeSecretSyncInterface interface {
	SecretSyncRules(namespace string) SecretSyncRuleInterface
}

type KubeSecretSyncClient struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*KubeSecretSyncClient, error) {
	AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &KubeSecretSyncClient{restClient: client}, nil
}

func (c *KubeSecretSyncClient) SecretSyncRules(namespace string) SecretSyncRuleInterface {
	return &secretSyncRuleClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
