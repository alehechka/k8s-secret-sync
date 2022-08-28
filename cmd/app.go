package cmd

import (
	kubesecretsync "github.com/alehechka/kube-secret-sync"

	"github.com/urfave/cli/v2"
)

// App represents the CLI application
func App() *cli.App {
	app := cli.NewApp()
	app.Version = kubesecretsync.Version
	app.Usage = "Automatically synchronize k8s Secrets across namespaces."
	app.Commands = []*cli.Command{
		StartCommand,
	}

	return app
}
