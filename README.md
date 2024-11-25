# Vault-Autounseal

A special application for checking the unsealed of the Vault

## Configuration

```yaml
---
checkInterval: 15
tlsSkipVerify: true
printUnsealLogs: false

nodes:
  - https://1.2.3.4:8200
  - https://1.2.3.5:8200
  - https://1.2.3.6:8200

unsealTokens:
  - token1
  - token2
  - token3
  - token4
  - token5
```

* `checkInterval` - the node verification interval in seconds.
* `tlsSkipVerify` - certificate verification is required when connecting.
* `printUnsealLogs` - print logs that the nodes have been printed.
* `nodes` - a list of nodes that need to be checked for unseal.
* `unsealTokens` - a list of unseal tokens.

## Deploy Vault-Autounseal

Using your vault or creating vault using [helm](https://developer.hashicorp.com/vault/docs/platform/k8s/helm/examples/ha-with-raft):
```bash
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update
helm upgrade --install -n vault --create-namespace vault hashicorp/vault \
    --set='server.ha.enabled=true' \
    --set='server.ha.raft.enabled=true'
kubectl exec -ti vault-0 -- vault operator init
kubectl exec -ti vault-0 -- vault operator unseal
kubectl -n vault exec -it vault-1 -- vault operator raft join http://vault-0.vault-internal:8200
kubectl -n vault exec -it vault-2 -- vault operator raft join http://vault-0.vault-internal:8200
```

Change unseal tokens in config:
```bash
vim deploy/vault-autounseal-config.yaml
```

Apply vault-autounseal:
```bash
kubectl -n vault apply -f deploy
```
