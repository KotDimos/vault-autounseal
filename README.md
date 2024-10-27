# Vault-Autounseal

Специальное приложение для автоматической распаковки хранилища Vault.

## Пример конфигурации

```yaml
---
checkInterval: 15
tlsSkipVerify: false
printUnsealLogs: true

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

* `checkInterval` - интервал проверки нод в секундах. 
* `tlsSkipVerify` - требуется ли проверка сертификатов при подключении.
* `printUnsealLogs` - печатать логи об том что ноды разпечатаны.
* `nodes` - список нод, которые нужно проверять на unseal.
* `unsealTokens` - список unseal токенов.
