# Tables for the database
## coin
|name|type|description|
|---|:--:|---|
|symbol|varchar(32)|coin's name, upper, primary key|
|decimal|tinyint|the coin's decimal|
|contract|varchar(256)|the address of contract, ***optional***|
|wallet|varchar(256)|address of coin's wallet|

## swap_pair
|name|type|description|
|---|:--:|---|
|id|int|increment, primary key|
|coin1|varchar(32)|table coin's 'symbol'|
|coin2|varchar(32)|table coin's 'symbol', coin1 and coin2 are composite keys|
|amount1|varchar(256)|coin1's amount, int256|
|amount2|varchar(256)|coin2's amount, int256|
|lp|varchar(256)|amount of lp token, int256|
|fee_rate|float|fee rate|

## balance
|name|type|description|
|---|:--:|---|
|account|varchar(256)|main account address, that user has private key, primary key|
|coin|varchar(32)|the coin name, table coin's symbol, primary key|
|amount|varchar(256)|amount of lp token, int256|
|ts|bigint|the last timestamp for update|

## pending_swap_order
|name|type|description|
|---|:--:|---|
|from_address|varchar(256)|user's account address|
|from_coin|varchar(32)|source coin name|
|from_amount|varchar(256)|amount of from_coin, int256|
|to_address|varchar(256)|target address of to_coin|
|to_coin|varchar(32)|target coin name|
|min_amount|varchar(256)|amount of to_coin, int256|
|o_time|bigint|order time for ms|

## swap_order
|name|type|description|
|---|:--:|---|
|id|bigint|increment, primary key|
|from_coin|varchar(32)|source coin name|
|from_amount|varchar(256)|amount of from_coin, int256|
|from_address|varchar(256)|user's account address|
|to_coin|varchar(32)|target coin name|
|to_amount|varchar(256)|amount of to_coin, int256|
|to_address|varchar(256)|target address of to_coin|
|fee|varchar(256)|fee, int256|
|state|tinyint|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 is the end state|
|o_time|bigint|order time for ms|
|e_time|bigint|end time for ms|

## pending_collect_order
|name|type|description|
|---|:--:|---|
|id|varchar(256)|hash, primary key|
|account|varchar(256)|user's account address|
|from|varchar(256)|the address of coin from|
|coin|varchar(32)|coin name|
|amount|varchar(256)|amount of coin to collect|
|state|tinyint|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 is the end state|
|o_time|bigint|order time for ms|
|e_time|bigint|end time for ms|

## collect_order
|name|type|description|
|---|:--:|---|
|id|varchar(256)|hash, primary key|
|account|varchar(256)|user's account address|
|from|varchar(256)|the address of coin from|
|coin|varchar(32)|coin name|
|amount|varchar(256)|amount of coin to collect|
|state|tinyint|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 is the end state|
|o_time|bigint|order time for ms|
|e_time|bigint|end time for ms|

## pending_retrieve_order
|name|type|description|
|---|:--:|---|
|id|varchar(256)|hash, primary key|
|account|varchar(256)|user's account address|
|to|varchar(256)|the address of coin to|
|coin|varchar(32)|coin name|
|amount|varchar(256)|amount of coin to collect|
|state|tinyint|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 is the end state|
|o_time|bigint|order time for ms|
|e_time|bigint|end time for ms|

## retrieve_order
|name|type|description|
|---|:--:|---|
|id|varchar(256)|hash, primary key|
|account|varchar(256)|user's account address|
|to|varchar(256)|the address of coin to|
|coin|varchar(32)|coin name|
|amount|varchar(256)|amount of coin to collect|
|state|tinyint|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 is the end state|
|o_time|bigint|order time for ms|
|e_time|bigint|end time for ms|

## pending_add_order
|name|type|description|
|---|:--:|---|
|id|varchar(256)|hash, primary key|
|account|varchar(256)|user's account address|
|coin1|string|name of coin1, default IOTA|
|coin2|string|name of coin2|
|amount|string|amount of IOAT|
|state|tinyint|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 is the end state|
|o_time|bigint|order time for ms|
|e_time|bigint|end time for ms|

## add_order
|name|type|description|
|---|:--:|---|
|id|varchar(256)|hash, primary key|
|account|varchar(256)|user's account address|
|coin1|string|name of coin1, default IOTA|
|coin2|string|name of coin2|
|amount|string|amount of IOAT|
|state|tinyint|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 is the end state|
|o_time|bigint|order time for ms|
|e_time|bigint|end time for ms|

## pending_remove_order
|name|type|description|
|---|:--:|---|
|id|varchar(256)|hash, primary key|
|account|varchar(256)|user's account address|
|coin1|string|name of coin1, default IOTA|
|coin2|string|name of coin2|
|amount|string|amount of lp token|
|state|tinyint|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 is the end state|
|o_time|bigint|order time for ms|
|e_time|bigint|end time for ms|

## remove_order
|name|type|description|
|---|:--:|---|
|id|varchar(256)|hash, primary key|
|account|varchar(256)|user's account address|
|coin1|string|name of coin1, default IOTA|
|coin2|string|name of coin2|
|amount|string|amount of lp token|
|state|tinyint|0:pending, 1:finished, 2:backing, 3:failed, 4:cancel. 1, 3 and 4 is the end state|
|o_time|bigint|order time for ms|
|e_time|bigint|end time for ms|