package client_test

import (
	"testing"

	pkg "github.com/alehechka/kube-secret-sync/client"
	"github.com/stretchr/testify/assert"
)

func Test_ListSecrets(t *testing.T) {
	client := InitializeTestClientset()

	client.CreateSecret(defaultNamespace, defaultSecret)

	secrets, err := client.ListSecrets(keyDefault)

	assert.NoError(t, err)
	assert.NotNil(t, secrets)
	assert.Equal(t, 1, len(secrets.Items))
}

func Test_ListSecrets_Empty(t *testing.T) {
	client := InitializeTestClientset()

	secrets, err := client.ListSecrets(keyDefault)

	assert.NoError(t, err)
	assert.NotNil(t, secrets)
	assert.Equal(t, 0, len(secrets.Items))
}

func Test_GetSecret(t *testing.T) {
	client := InitializeTestClientset()

	client.CreateSecret(defaultNamespace, defaultSecret)

	secret, err := client.GetSecret(keyDefault, keyDefaultSecret)

	assert.NoError(t, err)
	assert.NotNil(t, secret)
}

func Test_GetSecret_NoSecret(t *testing.T) {
	client := InitializeTestClientset()

	secret, err := client.GetSecret(keyDefault, keyDefaultSecret)

	assert.Error(t, err)
	assert.Nil(t, secret)
}

func Test_DeleteSecret(t *testing.T) {
	client := InitializeTestClientset()

	client.CreateSecret(defaultNamespace, defaultSecret)

	err := client.DeleteSecret(defaultNamespace, defaultSecret)

	assert.NoError(t, err)
}

func Test_DeleteSecret_NoSecret(t *testing.T) {
	client := InitializeTestClientset()

	err := client.DeleteSecret(defaultNamespace, defaultSecret)

	assert.Error(t, err)
}

func Test_UpdateSecret(t *testing.T) {
	client := InitializeTestClientset()

	client.CreateSecret(defaultNamespace, defaultSecret)

	err := client.UpdateSecret(defaultNamespace, defaultSecret)

	assert.NoError(t, err)
}

func Test_UpdateSecret_NoSecret(t *testing.T) {
	client := InitializeTestClientset()

	err := client.UpdateSecret(testNamespace, testSecret)

	assert.Error(t, err)
}

func Test_CreateSecret(t *testing.T) {
	client := InitializeTestClientset()

	err := client.CreateSecret(defaultNamespace, defaultSecret)
	assert.NoError(t, err)

	secret, err := client.GetSecret(keyDefault, keyDefaultSecret)
	assert.NoError(t, err)
	assert.NotNil(t, secret)
	assert.True(t, pkg.SecretsAreEqual(defaultSecret, secret))
}

func Test_IsManagedBy_True(t *testing.T) {
	secret := *defaultSecret
	secret.Annotations = managedByAnnotations
	isManaged := pkg.IsManagedBy(&secret)
	assert.True(t, isManaged)
}

func Test_IsManagedBy_False(t *testing.T) {
	isManaged := pkg.IsManagedBy(defaultSecret)
	assert.False(t, isManaged)
}
