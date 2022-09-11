package client_test

import (
	"testing"

	pkg "github.com/alehechka/kube-secret-sync/client"
	"github.com/alehechka/kube-secret-sync/constants"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
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

func Test_SecretsAreEqual_True(t *testing.T) {
	equal := pkg.SecretsAreEqual(defaultSecret, testSecret)
	assert.True(t, equal)
}

func Test_SecretsAreEqual_DifferentTypes(t *testing.T) {
	secret := *testSecret
	secret.Type = v1.SecretTypeOpaque

	equal := pkg.SecretsAreEqual(defaultSecret, &secret)
	assert.False(t, equal)
}

func Test_SecretsAreEqual_DifferentData(t *testing.T) {
	secret := *testSecret
	secret.Data = map[string][]byte{"item": {65, 66}}

	equal := pkg.SecretsAreEqual(defaultSecret, &secret)
	assert.False(t, equal)
}

func Test_SecretsAreEqual_DifferentStringData(t *testing.T) {
	secret := *testSecret
	secret.StringData = map[string]string{"something": "else"}

	equal := pkg.SecretsAreEqual(defaultSecret, &secret)
	assert.False(t, equal)
}

func Test_SecretsAreEqual_DifferentAnnotations(t *testing.T) {
	secret := *testSecret
	secret.Annotations = map[string]string{"something": "else"}

	equal := pkg.SecretsAreEqual(defaultSecret, &secret)
	assert.False(t, equal)
}

func Test_AnnotationsAreEqual_True(t *testing.T) {
	a := map[string]string{
		constants.ManagedByAnnotationKey:                constants.ManagedByAnnotationValue,
		constants.LastAppliedConfigurationAnnotationKey: "something-random",
		"same-key": "same-value",
	}

	b := map[string]string{
		constants.ManagedByAnnotationKey:                constants.ManagedByAnnotationValue,
		constants.LastAppliedConfigurationAnnotationKey: "something-different",
		"same-key": "same-value",
	}

	equal := pkg.AnnotationsAreEqual(a, b)
	assert.True(t, equal)
}

func Test_AnnotationsAreEqual_False(t *testing.T) {
	a := map[string]string{
		"same-key": "same-value",
	}

	b := map[string]string{
		"same-key": "different-value",
	}

	equal := pkg.AnnotationsAreEqual(a, b)
	assert.False(t, equal)
}

func Test_PrepareSecret(t *testing.T) {
	secret := *defaultSecret
	secret.Annotations = map[string]string{constants.LastAppliedConfigurationAnnotationKey: "some-previous-config", "some-key": "some-value"}

	prepared := pkg.PrepareSecret(testNamespace, &secret)

	assert.True(t, pkg.SecretsAreEqual(&secret, prepared))
	assert.True(t, pkg.IsManagedBy(prepared))

	lastConfig, ok := prepared.Annotations[constants.LastAppliedConfigurationAnnotationKey]
	assert.False(t, ok)
	assert.Equal(t, "", lastConfig)
}

func Test_CopyAnnotations(t *testing.T) {
	a := map[string]string{
		constants.ManagedByAnnotationKey:                constants.ManagedByAnnotationValue,
		constants.LastAppliedConfigurationAnnotationKey: "something-random",
		"same-key": "same-value",
	}

	copy := pkg.CopyAnnotations(a)

	assert.Equal(t, 1, len(copy))
	assert.Equal(t, "same-value", copy["same-key"])
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
