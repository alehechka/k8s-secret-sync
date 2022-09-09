package client_test

import (
	"context"
	"time"

	"github.com/alehechka/kube-secret-sync/client"
	"k8s.io/client-go/kubernetes/fake"
)

func InitializeTestClientset() *client.Client {
	c := new(client.Client)

	c.Context = context.Background()
	c.StartTime = time.Now()
	c.DefaultClientset = fake.NewSimpleClientset()

	return c
}
