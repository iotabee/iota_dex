# api address
```
    https://api.iotabee.com
```

# public api

## GET /public/pairs
### respose
```json
{
    "result":true,
    "data":[{
        "coin1":"IOTA",
        "coin2":"SMR",
        "reserve1":"10000000",
        "reserve1":"20000000",
        "total_supply":"102000",
        "fee_rate":0.003,
        "fee_scale":0.15
    }]
}
```

## GET /public/pair?coin1={1}&coin2={2}
### respose
```json
{
    "result":true,
    "data":{
        "coin1":"IOTA",
        "coin2":"SMR",
        "reserve1":"10000000",
        "reserve1":"20000000",
        "total_supply":"102000",
        "fee_rate":0.003,
        "fee_scale":0.15
    }
}
```

## GET /public/coin?symbol={}
### respose
```json
{
    "result":true,
    "data":{
        "contract":"contract address",
        "wallet":"wallet address",
        "decimal":18,
        "amount":"1000"
    }
}
```

# sign api

## A example for sign data with the private key
```go
    privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	data := []byte("1655714635") //
	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	sign := hexutil.Encode(signature)
	//0x930b692f4b3117d4f7e5640b6d19b383f29046ef6ffd38fe0c221065ab90c50e402037b99577f3469af5e1d507b3b9a00a23b3f2ee603c826a27cd30e58e9c9a00
```
Every private api must add the ts and sign params.
```
ts=1655714635&sign=0x930b692f4b3117d4f7e5640b6d19b383f29046ef6ffd38fe0c221065ab90c50e402037b99577f3469af5e1d507b3b9a00a23b3f2ee603c826a27cd30e58e9c9a00
```


## GET /balance
### respose
```json
{
    "result":true,
    "data":{
        "IOTA":"100",
        "SMR":"1000"
    }
}
```

## GET /order/swap
### request params
| name       |  type  | description                 |
| ---------- | :----: | --------------------------- |
| source     | string | source coin's name          |
| target     | string | target coin's name          |
| to         | string | address of 'target' to send |
| amount     | int256 | amount of 'source'          |
| min_amount | int256 | min amount of 'target'      |
### respose
```json
{
    "result":true
}
```

## GET /order/swap/pending get the swap order which is pending
### respose
```json
{
    "result":true,
    "data":{
        "from":"coin name",
        "from_amount":1,
        "from_address":"transfer in address",
        "to":"coin name",
        "to_amount":2,
        "to_address":"transfer out address",
        "fee":10,
        "o_time":1654418280
    }
}
```
| name  | type  | description                                                                       |
| ----- | :---: | --------------------------------------------------------------------------------- |
| state |  int  | 0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 id the end state |

## GET /order/swap/cancel cancel the swap order which is pending
### respose
```json
{
    "result":true
}
```

## GET /order/swap/list?count={10} get swap orders that had been dealed
### response
```json
{
    "result":true,
    "data":[{
        "id":1,
        "from":"coin name",
        "from_amount":1,
        "from_address":"transfer in address",
        "to":"coin name",
        "to_amount":2,
        "to_address":"transfer out address",
        "fee":10,
        "state":1,
        "o_time":1654418280,
        "e_time":1654418290
    }]
}
```

## GET /order/swap/info?id={1} get the swap order which had been dealed
### respose
```json
{
    "result":true,
    "data":{
        "id":1,
        "from":"coin name",
        "from_amount":1,
        "from_address":"transfer in address",
        "to":"coin name",
        "to_amount":2,
        "to_address":"transfer out address",
        "fee":10,
        "state":1,
        "o_time":1654418280,
        "e_time":1654418290
    }
}
```
| name  | type  | description                                                                       |
| ----- | :---: | --------------------------------------------------------------------------------- |
| state |  int  | 0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 id the end state |



## GET /order/coin/collect
### request params
| name    |  type  | description                         |
| ------- | :----: | ----------------------------------- |
| account | string | user's IOTA address                 |
| coin    | string | name of coin that will be collected |
| amount  | string | amount of collect                   |
### respose
```json
{
    "result":true
}
```

## GET /order/coin/retrieve
### request params
| name   |  type  | description                |
| ------ | :----: | -------------------------- |
| to     | string | to address                 |
| coin   | string | coin's name for retrieving |
| amount | string | amount for retrieving      |
### respose
```json
{
    "result":true,
    "data":1
}
```

## GET /order/coin/pending get the collect coin order pending
### respose
```json
{
    "result":true,
    "data":{
        "account":"the address of IOTA",
        "from":"from address",
        "coin":"",
        "amount":"",
        "o_time":1654418280,
    }
}
```

## GET /order/coin/cancel cancel the collect coin order which is pending
### respose
```json
{
    "result":true
}
```

## GET /order/coin/list?count={10} get coin orders that had been dealed
### response
```json
{
    "result":true,
    "data":[{
        "id":1,
        "account":"",
        "address":"",
        "coin":"",
        "amount":"",
        "direction":1,
        "state":1,
        "o_time":1654418280,
        "e_time":1654418290
    }]
}
```

## GET /order/coin/info?id={1} get the coin order which had been dealed
### respose
```json
{
    "result":true,
    "data":{
        "id":1,
        "account":"",
        "address":"",
        "coin":"",
        "amount":"",
        "direction":1,
        "state":1,
        "o_time":1654418280,
        "e_time":1654418290
    }
}
```

## GET /order/liquidity/add
### request params
| name    |  type  | description                 |
| ------- | :----: | --------------------------- |
| coin1   | string | name of coin1, default IOTA |
| coin2   | string | name of coin2               |
| amount1 | string | amount for coin1            |
### respose
```json
{
    "result":true
}
```

## GET /order/liquidity/remove
### request params
| name  |  type  | description                 |
| ----- | :----: | --------------------------- |
| coin1 | string | name of coin1, default IOTA |
| coin2 | string | name of coin2               |
| lp    | string | amount of lp token          |
### respose
```json
{
    "result":true,
    "data":1
}
```

## GET /order/liquidity/cancel cancel liquidity add order which is pending
### respose
```json
{
    "result":true,
}
```

## GET /order/liquidity/pending get the liquidity add order which is pending
### respose
```json
{
    "result":true,
    "data":{
        "account":"",
        "coin1":"",
        "coin2":"",
        "amount1":"",
        "o_time":1654418280
    }
}
```
| name    |  type  | description                                                          |
| ------- | :----: | -------------------------------------------------------------------- |
| amount1 | string | coin1's amount when direction=1, lp token's amount when direction=-1 |

## GET /order/liquidity/list?count={5} get liquidity orders that had been dealed
### respose
```json
{
    "result":true,
    "data":[{
        "id":0,
        "account":"",
        "coin1":"",
        "coin2":"",
        "amount1":"",
        "amount2":"",
        "lp":"",
        "direction":1,
        "state":1,
        "o_time":1654418280,
        "e_time":1654418280
    }]
}
```

## GET /order/liquidity/info?id={1} get liquidity order which had been dealed
### respose
```json
{
    "result":true,
    "data":{
        "id":1,
        "account":"",
        "coin1":"",
        "coin2":"",
        "amount1":"",
        "amount2":"",
        "lp":"",
        "direction":1,
        "state":1,
        "o_time":1654418280,
        "e_time":1654418280
    }
}
```

# Error Response
## A Error will return as a json string, when you require a api. For example.
```json
{
    "result":false,
    "err_code":1,
    "err_msg":"sign error"
}
```
| code | description  |
| ---- | :----------- |
| 1    | sige error   |
| 2    | params error |
| 3    | system error |
