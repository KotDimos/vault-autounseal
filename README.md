# Vault-Autounseal

A special application for checking the unsealed of the Vault

## Configuration

```yaml
---
checkInterval: 15
tlsSkipVerify: true
printUnsealLogs: fasle

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
* `printUnsealLogs` - печатать логи об том что ноды разпечатаны.
* `nodes` - a list of nodes that need to be checked for unseal.
* `unsealTokens` - a list of unseal tokens.

## Deploy Vault-Autounseal


