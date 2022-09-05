# kube-secret-sync

A kubernetes client-go app for syncing secrets across namespaces.

## Deployment

Examples deployment files can be found under [k8s](/k8s) and should be created in your cluster in the prefixed numerical order.

By using the example yamls out-of-the-box, a new `kube-secret-sync` namespace will be created and all secrets created within that namespace will be synced to all but the excluded namespaces (`kube-system`, `kube-public`, `kube-node-lease`). Additionally, the `DEBUG` flag is set to `true` which could be removed to limit verbosity of log messages.

## Configuration Options

| Environment Variable | Example            | Type      | Default | Description                                                                       |
| -------------------- | ------------------ | --------- | ------- | --------------------------------------------------------------------------------- |
| `POD_NAMESPACE`      | `custom-namespace` | `string`  |         | Specifies the namespace that current application pod is running in.               |
| `DEBUG`              | `true`             | `boolean` | `false` | Log debug messages.                                                               |
| `FORCE`              | `true`             | `boolean` | `false` | Forces synchronization of all secrets, not just kube-secret-sync managed secrets. |
