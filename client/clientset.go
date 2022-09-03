package client

import (
	"github.com/alehechka/kube-secret-sync/api/types/v1/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	DefaultClientset        *kubernetes.Clientset
	KubeSecretSyncClientset *clientset.KubeSecretSyncClient
)

func InitializeDefaultClientset(config *SyncConfig) error {
	cluster, err := clusterConfig(config)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(cluster)
	if err != nil {
		return err
	}

	DefaultClientset = clientset
	return nil
}

func InitializeKubeSecretSyncClient(config *SyncConfig) error {
	cluster, err := clusterConfig(config)
	if err != nil {
		return err
	}

	clientset, err := clientset.NewForConfig(cluster)
	if err != nil {
		return err
	}

	KubeSecretSyncClientset = clientset
	return nil
}

func clusterConfig(config *SyncConfig) (cluster *rest.Config, err error) {
	if config.OutOfCluster {
		cluster, err = clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	} else {
		cluster, err = rest.InClusterConfig()
	}

	return
}
