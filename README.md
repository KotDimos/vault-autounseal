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

* `checkInterval` - the node verification interval in seconds (default vaule - 15).
* `tlsSkipVerify` - certificate verification is required when connecting (default vaule - true).
* `printUnsealLogs` - if true, print logs that nodes are unsealed (default vaule - false).
* `nodes` - a list of nodes that need to be checked for unseal.
* `unsealTokens` - a list of unseal tokens.

## Deploy Vault-Autounseal

If you are using your vault, you are skipping this step.

Here is an example of deploying vault on integrated storage, if you want to use another option, look [here](https://developer.hashicorp.com/vault/docs/configuration/storage).

Creating vault using [helm](https://developer.hashicorp.com/vault/docs/platform/k8s/helm/examples/ha-with-raft):
```bash
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update
helm upgrade --install vault hashicorp/vault --create-namespace -n vault \
    --set='server.ha.enabled=true' \
    --set='server.ha.raft.enabled=true'
```

If helm hashicorp repo unavailable, cloning this [repo](https://github.com/hashicorp/vault-helm) and deploying vault:
```bash
git clone git@github.com:hashicorp/vault-helm.git
helm upgrade --install vault vault-helm --create-namespace -n vault \
    --set='server.ha.enabled=true' \
    --set='server.ha.raft.enabled=true'
```

Initialize vault:
```bash
kubectl -n vault exec -it vault-0 -- vault operator init
```

Unseal vault instance:
```bash
kubectl -n vault exec -it vault-0 -- vault operator unseal
kubectl -n vault exec -it vault-0 -- vault operator unseal
kubectl -n vault exec -it vault-0 -- vault operator unseal
```

Join:
```bash
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
