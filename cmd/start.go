package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

const (
	debugFlag             = "debug"
	excludeSecretsFlag    = "exclude-secrets"
	includeSecretsFlag    = "include-secrets"
	excludeNamespacesFlag = "exclude-namespaces"
	includeNamespacesFlag = "include-namespaces"
	secretsNamespaceFlag  = "secrets-namespace"
)

var startFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    debugFlag,
		Usage:   "Log debug messages.",
		EnvVars: []string{"DEBUG"},
	},
	&cli.StringSliceFlag{
		Name:    excludeSecretsFlag,
		Usage:   "Excludes specific Secrets from syncing. Will override `included` Secrets if specified in both. Supply as CSV in environment variables.",
		EnvVars: []string{"EXCLUDE_SECRETS"},
	},
	&cli.StringSliceFlag{
		Name:    includeSecretsFlag,
		Usage:   "Includes specific Secrets in syncing. Acts as a whitelist and all other Secrets will not be synced. Supply as CSV in environment variables.",
		EnvVars: []string{"INCLUDE_SECRETS"},
	},
	&cli.StringSliceFlag{
		Name:    excludeNamespacesFlag,
		Usage:   "Excludes specific Namespaces from syncing. Will override `included` Namespaces if specified in both. Supply as CSV in environment variables.",
		EnvVars: []string{"EXCLUDE_NAMESPACES"},
	},
	&cli.StringSliceFlag{
		Name:    includeNamespacesFlag,
		Usage:   "Includes specific Namespaces in syncing. Acts as a whitelist and all other Namespaces will not be synced. Supply as CSV in environment variables.",
		EnvVars: []string{"INCLUDE_NAMESPACES"},
	},
	&cli.StringFlag{
		Name:    secretsNamespaceFlag,
		Usage:   "Specifies which namespace to sync secrets from, defaults to the same namespace that this application is deployed in.",
		EnvVars: []string{"SECRETS_NAMESPACE"},
	},
}

func startKubeSecretSync(ctx *cli.Context) (err error) {
	if ctx.Bool(debugFlag) {
		log.SetLevel(log.DebugLevel)
	}

	return nil
}

// StartCommand starts the kube-secret-sync process.
var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "Start kube-secret-sync application.",
	Action: startKubeSecretSync,
	Flags:  startFlags,
}
