# public api

## GET /pairs
### respose
```json
{
    "result":true,
    "data":[{
        "coins":{
            "IOTA":{
                "contract":"",
                "wallet":"",
                "amount":"5230498034582345",
                "decimal":18
            },
            "SMR":{
                "contract":"",
                "wallet":"",
                "amount":"5230498034582345",
                "decimal":18
            }
        },
        "lp":"102000",
        "fee_rate":0.003
    }]
}
```

## GET /price?coin1={1}&coin2={2}
### respose
```json
{
    "result":true,
    "coin1":123,
    "coin2":124,
}
```

## GET /balance?account={1}
### respose
```json
{
    "result":true,
    "data":{
        "IOTA":100,
        "SMR":1000
    }
}
```

## GET /market/swap/list?count={10} 
### request params
```
max count is 100
```
### respose
```json
{
    "result":true,
    "data":[{
        "id":1,
        "from_address":"",
        "from_coin":"coin name",
        "from_amount":1,
        "to_address":"",
        "to_coin":"coin name",
        "to_amount":2,
        "end_time":1654418290,
    }]
}
```

# sign api

## A example for sign data with the private key
```go
    privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	data := []byte("19152") //total days since 1970.01.01 00:00:00
	hash := crypto.Keccak256Hash(data) //0xb2572721d4d9b79d1f7678b68da3265dfc42dbae0ce0603cc17e3af26f5deb51
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	sign := hexutil.Encode(signature)
	//0x76007ce42c2cfd60f5689ea17179ea8cf3df8050c582fe0a0df9de36afcba4026cfe2bfe09012629db451a8abcddfa3fb0fc327c37402c73c6b133526a19b2fc01
```

## GET /order/swap
### request params
|name|type|description|
|---|:--:|---|
|source|string|source coin's name|
|target|string|target coin's name|
|to|string|address of 'target' to send|
|amount|int256|amount of 'source'|
|min_amount|int256|min amount of 'target'|
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
|name|type|description|
|---|:--:|---|
|state|int|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 id the end state|

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
|name|type|description|
|---|:--:|---|
|state|int|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 id the end state|



## GET /order/coin/collect
### request params
|name|type|description|
|---|:--:|---|
|account|string|user's IOTA address|
|coin|string|name of coin that will be collected|
|amount|string|amount of collect|
### respose
```json
{
    "result":true
}
```

## GET /order/coin/retrieve
### request params
|name|type|description|
|---|:--:|---|
|to|string|to address|
|coin|string|coin's name for retrieving|
|amount|string|amount for retrieving|
### respose
```json
{
    "result":true,
    "id":"1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8"
}
```



## GET /order/coin/pending get the coin order which is pending
### respose
```json
{
    "result":true,
    "data":{
        "account":"",
        "address":"",
        "coin":"",
        "amount":"",
        "direction":1,
        "o_time":1654418280,
    }
}
```
|name|type|description|
|---|:--:|---|
|account|string|the address of IOTA|
|address|string|from address when direction=1 and to address when direction=-1|
|direction|int|1: collect, -1: retrieve|

## GET /order/coin/cancel cancel the coin order which is pending
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
|name|type|description|
|---|:--:|---|
|coin1|string|name of coin1, default IOTA|
|coin2|string|name of coin2|
|amount|string|amount for coin1|
### respose
```json
{
    "result":true
}
```

## GET /order/liquidity/remove
### request params
|name|type|description|
|---|:--:|---|
|coin1|string|name of coin1, default IOTA|
|coin2|string|name of coin2|
|lp|string|amount of lp token|
### respose
```json
{
    "result":true,
    "id":"1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8"
}
```

## GET /order/liquidity/cancel cancel liquidity order which is pending
### request params
|name|type|description|
|---|:--:|---|
### respose
```json
{
    "result":true,
}
```

## GET /order/liquidity/pending get the liquidity order which is pending
### respose
```json
{
    "result":true,
    "data":{
        "account":"",
        "coin1":"",
        "coin2":"",
        "amount1":"",
        "direction":1,
        "o_time":1654418280
    }
}
```
|name|type|description|
|---|:--:|---|
|amount1|string|coin1's amount when direction=1, lp token's amount when direction=-1|

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
|code|description|
|---|:---|
|1|sige error|
|2|params error|
|3|system error|

{
	"env":"dev",
	"http_port":8080,
	"db": {
		"lifetime": 30,
		"openconns": 200,
		"idleconns": 100,
		"host": "",
		"port": "3306",
		"dbname": "iota_dex",
		"usr": "root",
		"pwd": "851012"
	},
	"token_time":7200
}
