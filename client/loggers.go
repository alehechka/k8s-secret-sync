package client

import (
	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func ruleLogger(rule *typesv1.SecretSyncRule) *log.Entry {
	return ruleNameLogger(rule.Name)
}

func ruleNameLogger(name string) *log.Entry {
	return log.WithFields(log.Fields{"name": name, "kind": "SecretSyncRule"})
}

func namespaceLogger(namespace *v1.Namespace) *log.Entry {
	return namespaceNameLogger(namespace.Name)
}

func namespaceNameLogger(namespace string) *log.Entry {
	return log.WithFields(log.Fields{"name": namespace, "kind": "Namespace"})
}

func secretLogger(secret *v1.Secret, namespaces ...*v1.Namespace) *log.Entry {
	namespace := secret.Namespace
	if len(namespaces) > 0 {
		namespace = namespaces[0].Name
	}

	return secretNameLogger(namespace, secret.Name)
}

func secretNameLogger(namespace, name string) *log.Entry {
	return log.WithFields(log.Fields{"name": name, "kind": "Secret", "namespace": namespace})
}
