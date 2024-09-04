# APIs

- [APIs](#apis)
  - [NewClient](#newclient)
  - [Hardwares](#hardwares)

## NewClient

Definition:
Creates a Swan Client instance

```shell
func NewClient(apiKey, isTestnet) *APIClient
```
Inputs:

| name      | type   | description                                   |
| --------- | ------ | --------------------------------------------- |
| apiKey    | string | Swan API key                                  |
| isTestnet | bool   | If set to true use testnet, otherwise maninet |

Outputs:

```shell
*APIClient            # Created swan Client instance.
```

## Hardwares

`Hardwares` Fetch available instance resources

```go
func (c *APIClient) Hardwares() ([]*Hardware, error) 
```

Inputs:

| name      | type   | description                                   |
| --------- | ------ | --------------------------------------------- |
| apiKey    | string | Swan API key                                  |
| isTestnet | bool   | If set to true use testnet, otherwise maninet |

Outputs:

| name      | type   | description                                   |
| --------- | ------ | --------------------------------------------- |
| apiKey    | string | Swan API key                                  |
| isTestnet | bool   | If set to true use testnet, otherwise maninet |

