package client_test

import (
	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	"github.com/alehechka/kube-secret-sync/constants"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	keyTestNamespace string = "test-namespace"
	keyTestSecret    string = "test-secret"
	keyDefault       string = "default"
	keyDefaultSecret string = "default-secret"
)

var managedByAnnotations = map[string]string{constants.ManagedByAnnotationKey: constants.ManagedByAnnotationValue}

var testNamespace = &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: keyTestNamespace}}
var defaultNamespace = &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: keyDefault}}
var defaultSecret = &v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: keyDefaultSecret, Namespace: keyDefault}}
var testSecret = &v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: keyTestSecret, Namespace: keyTestSecret}}
var testSecretSyncRule = &typesv1.SecretSyncRule{Spec: typesv1.SecretSyncRuleSpec{Secret: typesv1.Secret{Name: keyDefaultSecret, Namespace: keyDefault}}}
