# kube-secret-sync

A kubernetes client-go app for syncing secrets across namespaces.

## Deployment

Examples deployment files can be found under [k8s](/k8s) and should be created in your cluster in the prefixed numerical order.

By using the example yamls out-of-the-box, a new `kube-secret-sync` namespace will be created and all secrets created within that namespace will be synced to all but the excluded namespaces (`kube-system`, `kube-public`, `kube-node-lease`). Additionally, the `DEBUG` flag is set to `true` which could be removed to limit verbosity of log messages.

## Configuration Options

| Environment Variable | Example | Type | Default | Description |
| -------------------- | ------- | ---- | ------- | ----------- |
| `SECRETS_NAMESPACE`  | `custom-secret-namespace` | `string` | `default` | Specifies which namespace to sync secrets from. |
| `EXCLUDE_SECRETS`    | `do-not-sync, super-secret` | `string(csv)` | | Excludes specific Secrets from syncing. Will override **included** Secrets if specified in both. |
| `INCLUDE_SECRETS`    | `syncable-secret, other-secret` | `string(csv)` | | Includes specific Secrets in syncing. Acts as a whitelist and all other Secrets will not be synced. |
| `EXCLUDE_NAMESPACES` | `kube-system, kube-public` | `string(csv)` | | Excludes specific Namespaces from syncing. Will override **included** Namespaces if specified in both. |
| `INCLUDE_NAMESPACES` | `default, my-namespace` | `string(csv)` | | Includes specific Namespaces in syncing. Acts as a whitelist and all other Namespaces will not be synced. |
| `DEBUG` | `true` | `boolean` | `false` | Log debug messages. |
| `FORCE` | `true` | `boolean` | `false` | Forces synchronization of all secrets, not just kube-secret-sync managed secrets. |
