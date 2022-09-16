package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_SyncSecretToNamespace(t *testing.T) {
	client := InitializeTestClientset()

	client.CreateSecret(testSecretSyncRule, defaultNamespace, defaultSecret)

	err := client.SyncSecretToNamespace(testSecretSyncRule, testNamespace)

	assert.NoError(t, err)
}

func Test_SyncSecretToNamespace_NoSecret(t *testing.T) {
	client := InitializeTestClientset()

	err := client.SyncSecretToNamespace(testSecretSyncRule, testNamespace)

	assert.Error(t, err)
}

func Test_ListNamespaces(t *testing.T) {
	client := InitializeTestClientset()

	client.DefaultClientset.CoreV1().Namespaces().Create(client.Context, testNamespace, metav1.CreateOptions{})

	namespaces, err := client.ListNamespaces()
	assert.NoError(t, err)
	assert.NotNil(t, namespaces)
	assert.Equal(t, 1, len(namespaces.Items))
}

func Test_ListNamespaces_Empty(t *testing.T) {
	client := InitializeTestClientset()

	namespaces, err := client.ListNamespaces()
	assert.NoError(t, err)
	assert.NotNil(t, namespaces)
	assert.Equal(t, 0, len(namespaces.Items))
}
