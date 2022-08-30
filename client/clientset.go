package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func clientset(config *SyncConfig) (*kubernetes.Clientset, error) {
	cluster, err := clusterConfig(config)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(cluster)
}

func clusterConfig(config *SyncConfig) (cluster *rest.Config, err error) {
	if config.OutOfCluster {
		cluster, err = clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	} else {
		cluster, err = rest.InClusterConfig()
	}

	return
}
