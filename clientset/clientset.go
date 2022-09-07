package clientset

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	Default *kubernetes.Clientset
)

func InitializeDefault(cluster *rest.Config) error {
	clientset, err := kubernetes.NewForConfig(cluster)
	if err != nil {
		return err
	}

	Default = clientset
	return nil
}

func ClusterConfig(config *SyncConfig) (cluster *rest.Config, err error) {
	if config.OutOfCluster {
		cluster, err = clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	} else {
		cluster, err = rest.InClusterConfig()
	}

	return
}
