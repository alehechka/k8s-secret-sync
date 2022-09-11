# kube-secret-sync

A kubernetes client-go app for syncing secrets across namespaces.

The application works by creating `SecretSyncRule` resources that define which secrets to sync and which namespaces to sync it to. The application will watch all `SecretSyncRule`, `Secret`, and `Namespace` resources for changes and sync the defined secrets where necessary.

## Deployment

Each released version will create a package helm chart and compiled yaml to deploy the entire application and its required resources. The most recent version can be deployed with the following:

### `helm`

```bash
helm install kube-secret-sync https://github.com/alehechka/kube-secret-sync/releases/download/v1.1.0/kube-secret-sync-1.1.0.tgz --namespace kube-secret-sync --create-namespace
```

### `kubectl`

```bash
kubectl apply -f https://github.com/alehechka/kube-secret-sync/releases/download/v1.1.0/kube-secret-sync.yaml
```

## Secret Sync Rules

Syncing secrets is an opt-in process per secret. This allows for fine-grained control of which secrets get synced and where to. The below example the `my-api-key` secret will be automatically synced from the `default` namespace to all but the `kube-[.]*` system defined namespaces:

```yaml
apiVersion: kube-secret-sync.io/v1
kind: SecretSyncRule
metadata:
  name: my-api-key-rule
spec:
  secret: my-api-key
  namespace: default
  rules:
    namespaces:
      excludeRegex:
        - 'kube-[.]*'
```

Full `SecretSyncRule` configuration options

| Spec Variable                   | Example            | Type       | Description                                                                                                                                                      |
| ------------------------------- | ------------------ | ---------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `secret`                        | `mysecret`         | `string`   | The name of a secret to sync.                                                                                                                                    |
| `namespace`                     | `default`          | `string`   | The name of the namespace that the secret to sync is defined in.                                                                                                 |
| `rules.namespaces.exclude`      | `["kube-system"]`  | `[]string` | A list of namespaces to exclude from syncing (will take precedence over include rules).                                                                          |
| `rules.namespaces.excludeRegex` | `["kube-[.]*"]`    | `[]string` | A list of regex patterns that represent namespaces to exclude from syncing (will take precedence over include rules).                                            |
| `rules.namespaces.include`      | `["my-namespace"]` | `[]string` | A list of namespaces to include in syncing (all non-included will be excluded).                                                                                  |
| `rules.namespaces.includeRegex` | `["my-[.]"]`       | `[]string` | A list of regex patterns that represent namespaces to include in syncing (all non-included will be excluded).                                                    |
| `rules.force`                   | `true`             | `boolean`  | A flag to turn on forced sync. By default, the app will only sync "managed-by" secrets. With this on, any non-managed matching secret names will also be synced. |

## Configuration Options

The application itself has a few configuration options, however these are mainly used during local development and should most likely not be changed.

| Environment Variable | Example            | Type      | Default            | Description                                                         |
| -------------------- | ------------------ | --------- | ------------------ | ------------------------------------------------------------------- |
| `POD_NAMESPACE`      | `custom-namespace` | `string`  | `kube-secret-sync` | Specifies the namespace that current application pod is running in. |
| `DEBUG`              | `true`             | `boolean` | `false`            | Log debug messages.                                                 |

# Contribute

- Create an issue or open a pull request

# Author

[Adam Lehechka](https://github.com/alehechka)

# License

[MIT](/LICENSE)
