package cmd

import (
	"path/filepath"

	"github.com/alehechka/kube-secret-sync/api/types"
	"github.com/alehechka/kube-secret-sync/client"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/homedir"
)

const (
	debugFlag                  = "debug"
	excludeSecretsFlag         = "exclude-secrets"
	excludeRegexSecretsFlag    = "exclude-regex-secrets"
	includeSecretsFlag         = "include-secrets"
	includeRegexSecretsFlag    = "include-regex-secrets"
	excludeNamespacesFlag      = "exclude-namespaces"
	excludeRegexNamespacesFlag = "exclude-regex-namespaces"
	includeNamespacesFlag      = "include-namespaces"
	includeRegexNamespacesFlag = "include-regex-namespaces"
	secretsNamespaceFlag       = "secrets-namespace"
	outOfClusterFlag           = "out-of-cluster"
	kubeconfigFlag             = "kubeconfig"
	forceSyncFlag              = "force"
)

func kubeconfig() *cli.StringFlag {
	kubeconfig := &cli.StringFlag{Name: kubeconfigFlag}
	if home := homedir.HomeDir(); home != "" {
		kubeconfig.Value = filepath.Join(home, ".kube", "config")
		kubeconfig.Usage = "(optional) absolute path to the kubeconfig file"
	} else {
		kubeconfig.Usage = "absolute path to the kubeconfig file (required if running OutOfCluster)"
	}
	return kubeconfig
}

var startFlags = []cli.Flag{
	kubeconfig(),
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
		Name:    excludeRegexSecretsFlag,
		Usage:   "Excludes specific Secrets from syncing using regex matching. Will override `included` Secrets if specified in both. Supply as CSV in environment variables.",
		EnvVars: []string{"EXCLUDE_REGEX_SECRETS"},
	},
	&cli.StringSliceFlag{
		Name:    includeSecretsFlag,
		Usage:   "Includes specific Secrets in syncing. Acts as a whitelist and all other Secrets will not be synced. Supply as CSV in environment variables.",
		EnvVars: []string{"INCLUDE_SECRETS"},
	},
	&cli.StringSliceFlag{
		Name:    includeRegexSecretsFlag,
		Usage:   "Includes specific Secrets in syncing using regex matching. Acts as a whitelist and all other Secrets will not be synced. Supply as CSV in environment variables.",
		EnvVars: []string{"INCLUDE_REGEX_SECRETS"},
	},
	&cli.StringSliceFlag{
		Name:    excludeNamespacesFlag,
		Usage:   "Excludes specific Namespaces from syncing. Will override `included` Namespaces if specified in both. Supply as CSV in environment variables.",
		EnvVars: []string{"EXCLUDE_NAMESPACES"},
	},
	&cli.StringSliceFlag{
		Name:    excludeRegexNamespacesFlag,
		Usage:   "Excludes specific Namespaces from syncing using regex matching. Will override `included` Namespaces if specified in both. Supply as CSV in environment variables.",
		EnvVars: []string{"EXCLUDE_REGEX_NAMESPACES"},
	},
	&cli.StringSliceFlag{
		Name:    includeNamespacesFlag,
		Usage:   "Includes specific Namespaces in syncing. Acts as a whitelist and all other Namespaces will not be synced. Supply as CSV in environment variables.",
		EnvVars: []string{"INCLUDE_NAMESPACES"},
	},
	&cli.StringSliceFlag{
		Name:    includeRegexNamespacesFlag,
		Usage:   "Includes specific Namespaces in syncing using regex matching. Acts as a whitelist and all other Namespaces will not be synced. Supply as CSV in environment variables.",
		EnvVars: []string{"INCLUDE_REGEX_NAMESPACES"},
	},
	&cli.StringFlag{
		Name:    secretsNamespaceFlag,
		Usage:   "Specifies which namespace to sync secrets from.",
		EnvVars: []string{"SECRETS_NAMESPACE"},
		Value:   v1.NamespaceDefault,
	},
	&cli.BoolFlag{
		Name:    outOfClusterFlag,
		Usage:   "Will use the default ~/.kube/config file on the local machine to connect to the cluster externally.",
		Aliases: []string{"local"},
	},
	&cli.BoolFlag{
		Name:    forceSyncFlag,
		Usage:   "Forces synchronization of all secrets, not just kube-secret-sync managed secrets.",
		EnvVars: []string{"FORCE"},
	},
}

func startKubeSecretSync(ctx *cli.Context) (err error) {
	if ctx.Bool(debugFlag) {
		log.SetLevel(log.DebugLevel)
	}

	excludeRegexSecrets, err := types.CompileAll(ctx.StringSlice(excludeRegexSecretsFlag))
	if err != nil {
		return err
	}

	includeRegexSecrets, err := types.CompileAll(ctx.StringSlice(includeRegexSecretsFlag))
	if err != nil {
		return err
	}

	excludeRegexNamespaces, err := types.CompileAll(ctx.StringSlice(excludeRegexNamespacesFlag))
	if err != nil {
		return err
	}

	includeRegexNamespaces, err := types.CompileAll(ctx.StringSlice(includeRegexNamespacesFlag))
	if err != nil {
		return err
	}

	return client.SyncSecrets(&client.SyncConfig{
		ExcludeSecrets:      ctx.StringSlice(excludeSecretsFlag),
		ExcludeRegexSecrets: excludeRegexSecrets,
		IncludeSecrets:      ctx.StringSlice(includeSecretsFlag),
		IncludeRegexSecrets: includeRegexSecrets,

		ExcludeNamespaces:      ctx.StringSlice(excludeNamespacesFlag),
		ExcludeRegexNamespaces: excludeRegexNamespaces,
		IncludeNamespaces:      ctx.StringSlice(includeNamespacesFlag),
		IncludeRegexNamespaces: includeRegexNamespaces,

		SecretsNamespace: ctx.String(secretsNamespaceFlag),

		ForceSync: ctx.Bool(forceSyncFlag),

		OutOfCluster: ctx.Bool(outOfClusterFlag),
		KubeConfig:   ctx.String(kubeconfigFlag),
	})
}

// StartCommand starts the kube-secret-sync process.
var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "Start the kube-secret-sync application.",
	Action: startKubeSecretSync,
	Flags:  startFlags,
}
