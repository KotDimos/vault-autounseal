# vault-autounseal dev

1) Start vault:
```bash
docker compose up
```

2) In logs find Unseal tokens
```
Unseal Key: USNEAL_TOKEN
````

And change `vault-unseal.yaml` config

```yaml
unsealTokens:
  - UNSEAL_TOKEN
```

3) Build golang app
```bash
go build .
```

4) Start vault-autounseal
```bash
./vault-autounseal --config dev/vault-unseal.yaml
```

5) Special sealed vault.
```bash
docker exec -it vault sh
VAULT_TOKEN=<ROOT_TOKEN> vault operator seal
```

6) Check the logs vault-autounseal that autounseal is running
```
Node 'http://localhost:8200' is seal, start unseal
```
