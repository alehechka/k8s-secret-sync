package client

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	kssclientset "github.com/alehechka/kube-secret-sync/api/types/v1/clientset"
	"github.com/alehechka/kube-secret-sync/clientset"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func initClientsets(config *clientset.SyncConfig) error {
	clusterConfig, err := clientset.ClusterConfig(config)
	if err != nil {
		return err
	}

	if err := clientset.InitializeDefault(clusterConfig); err != nil {
		return err
	}

	if err := kssclientset.InitializeKubeSecretSync(clusterConfig); err != nil {
		return err
	}

	return nil
}

func initWatchers(ctx context.Context) (secretWatcher watch.Interface, namespaceWatcher watch.Interface, secretSyncRuleWatcher watch.Interface, err error) {
	secretWatcher, err = clientset.Default.CoreV1().Secrets(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	namespaceWatcher, err = clientset.Default.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	secretSyncRuleWatcher, err = kssclientset.KubeSecretSync.SecretSyncRules().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	return
}

func initSignalChannel() (sigc chan os.Signal) {
	sigc = make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	return
}
