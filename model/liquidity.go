package model

import (
	"errors"
	"fmt"
	"iota_dex/config"
	"math/big"
	"time"
)

type PendingLiquidityAddOrder struct {
	Account   string `json:"account"`
	Coin1     string `json:"coin1"`
	Coin2     string `json:"coin2"`
	Amount1   string `json:"amount1"`
	OrderTime int    `json:"o_time"`
}

type LiquidityOrder struct {
	Id        int64  `json:"id"`
	Account   string `json:"account"`
	Coin1     string `json:"coin1"`
	Coin2     string `json:"coin2"`
	Amount1   string `json:"amount1"`
	Amount2   string `json:"amount2"`
	Lp        string `json:"lp"`
	Direction int    `json:"direction"`
	State     int    `json:"state"`
	OrderTime int    `json:"o_time"`
	EndTime   int    `json:"e_time"`
}

//AddLiquidity add liquidity to pool
func AddLiquidity(account, coin1, coin2 string, amount1 *big.Int) error {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return err
	}

	//query account's balances to get the amounts for coin1,coin2 and lp-coin
	rows, err := tx.Query("select coin,amount from balance where account=? for update", account)
	if err != nil {
		tx.Rollback()
		return err
	}
	balances := make(map[string]*big.Int)
	for rows.Next() {
		var str1, str2 string
		if err = rows.Scan(&str1, &str2); err != nil {
			tx.Rollback()
			return err
		}
		a, b := new(big.Int).SetString(str2, 10)
		if !b {
			tx.Rollback()
			return fmt.Errorf("the balance amount(%s : %s) convert To big.Int error", str1, str2)
		}
		balances[str1] = a
	}
	if b, exist := balances[coin1]; !exist || b.Cmp(amount1) < 0 {
		tx.Rollback()
		return errors.New("coin1's amount is not enough")
	}

	//query the reserve and lp from the trade pair pool
	c1, c2 := coin1, coin2
	if c1 > c2 {
		c1, c2 = c2, c1
	}
	row := tx.QueryRow("select reserve1,reserve2,total_supply from swap_pair where coin1=? and coin2=? for update", c1, c2)
	var str1, str2, str3 string
	if err = row.Scan(&str1, &str2, &str3); err != nil {
		tx.Rollback()
		return err
	}
	reserve1, b1 := new(big.Int).SetString(str1, 10)
	reserve2, b2 := new(big.Int).SetString(str2, 10)
	totalSupply, b3 := new(big.Int).SetString(str3, 10)
	if !b1 || !b2 || !b3 || totalSupply.Cmp(big.NewInt(0)) <= 0 {
		tx.Rollback()
		return fmt.Errorf("reserve1,reserve2,totalSupply convert to big.Int error. %s:%s:%s", str1, str2, str3)
	}
	if c1 != coin1 {
		reserve1, reserve2 = reserve2, reserve1
	}

	//caculate the amount of coin2 for need
	a2 := new(big.Int).Div(new(big.Int).Mul(amount1, reserve2), reserve1)
	if balances[coin2] == nil || balances[coin2].Cmp(a2) < 0 {
		tx.Rollback()
		return fmt.Errorf("coin2's amount is not enough (%v, %s)", balances[coin2], a2.String())
	}

	//caculate the lp amount to mint, and upate the pool's state
	addLiquidity := new(big.Int).Div(new(big.Int).Mul(amount1, totalSupply), reserve1)
	totalSupply.Add(totalSupply, addLiquidity)
	reserve1.Add(reserve1, amount1)
	reserve2.Add(reserve2, a2)
	if c1 != coin1 {
		reserve1, reserve2 = reserve2, reserve1
	}
	if _, err := tx.Exec("update swap_pair set reserve1=?,reserve2=?,total_supply=? where coin1=? and coin2=?", reserve1.String(), reserve2.String(), totalSupply.String(), c1, c2); err != nil {
		tx.Rollback()
		return err
	}

	//update the coin1's amount
	balances[coin1] = new(big.Int).Sub(balances[coin1], amount1)
	if _, err := tx.Exec("update balance set amount=? where account=? and coin=?", balances[coin1].String(), account, coin1); err != nil {
		tx.Rollback()
		return err
	}

	//update the coin2's amount
	balances[coin2] = new(big.Int).Sub(balances[coin2], a2)
	if _, err := tx.Exec("update balance set amount=? where account=? and coin=?", balances[coin2].String(), account, coin2); err != nil {
		tx.Rollback()
		return err
	}

	//update the lpToken's amount
	lpCoin := "LP-" + c1 + "-" + c2
	if lp, exist := balances[lpCoin]; exist {
		balances[lpCoin] = new(big.Int).Add(lp, addLiquidity)
	} else {
		balances[lpCoin] = addLiquidity
	}
	if _, err := tx.Exec("replace into balance(account,coin,amount) values(?,?,?)", account, lpCoin, balances[lpCoin].String()); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("insert into `liquidity_order`(`account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`) VALUES(?,?,?,?,?,?,-1,1,?)", account, coin1, coin2, amount1.String(), a2.String(), addLiquidity.String(), time.Now().Unix()); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

//InsertPendingLiquidityAddOrder add a record to liquidity_add_order_pending, waiting for deal
func InsertPendingLiquidityAddOrder(account, coin1, coin2, amount string) error {
	if _, err := db.Exec("insert into `liquidity_add_order_pending`(`account`,`coin1`,`coin2`,`amount1`) VALUES(?,?,?,?)", account, coin1, coin2, amount); err != nil {
		return err
	}
	return nil
}

//RemoveLiquidity remove liquidity from pool
//@amount	: amount of lp token to remove
//@coin1	: require coin1 < coin2
func RemoveLiquidity(account, coin1, coin2 string, amount *big.Int) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return 0, err
	}

	//lock table balance to get the amount of lp token
	lpCoin := "LP-" + coin1 + "-" + coin2
	rows, err := tx.Query("select coin,amount from balance where account=? and (coin=? or coin=? or coin=?)for update", account, coin1, coin2, lpCoin)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	balances := make(map[string]*big.Int)
	for rows.Next() {
		var str1, str2 string
		if err = rows.Scan(&str1, &str2); err != nil {
			tx.Rollback()
			return 0, err
		}
		a, b := new(big.Int).SetString(str2, 10)
		if !b {
			tx.Rollback()
			return 0, fmt.Errorf("the balance amount(%s : %s) convert To big.Int error", str1, str2)
		}
		balances[str1] = a
	}
	if b, exist := balances[lpCoin]; !exist || b.Cmp(amount) < 0 {
		tx.Rollback()
		return 0, errors.New("amount of lp token is not enough")
	}

	//lock table swap_pair to get the amount of reserve and lp token
	row := tx.QueryRow("select reserve1,reserve2,total_supply from swap_pair where coin1=? and coin2=? for update", coin1, coin2)
	var str1, str2, str3 string
	if err = row.Scan(&str1, &str2, &str3); err != nil {
		tx.Rollback()
		return 0, err
	}
	reserve1, b1 := new(big.Int).SetString(str1, 10)
	reserve2, b2 := new(big.Int).SetString(str2, 10)
	totalSupply, b3 := new(big.Int).SetString(str3, 10)
	if !b1 || !b2 || !b3 || totalSupply.Cmp(amount) <= 0 {
		tx.Rollback()
		return 0, fmt.Errorf("reserve1,reserve2,totalSupply convert to big.Int error. %s:%s:%s", str1, str2, str3)
	}

	//caculate the out amounts for coin1 and coin2
	amountOut1 := new(big.Int).Div(new(big.Int).Mul(amount, reserve1), totalSupply)
	amountOut2 := new(big.Int).Div(new(big.Int).Mul(amount, reserve2), totalSupply)
	if amountOut1.Cmp(big.NewInt(0)) <= 0 || amountOut2.Cmp(big.NewInt(0)) <= 0 {
		tx.Rollback()
		return 0, fmt.Errorf("insufficient liquidity burned")
	}

	//update the swap_pair's reserve and burn lp token
	if _, err := tx.Exec("update swap_pair set reserve1=?,reserve2=?,total_supply=? where coin1=? and coin2=?", reserve1.Sub(reserve1, amountOut1).String(), reserve2.Sub(reserve2, amountOut2).String(), totalSupply.Sub(totalSupply, amount).String(), coin1, coin2); err != nil {
		tx.Rollback()
		return 0, err
	}

	//update the amount of lp token
	lp := new(big.Int).Sub(balances[lpCoin], amount)
	if _, err := tx.Exec("update balance set amount=? where account=? and coin=?", lp.String(), account, lpCoin); err != nil {
		tx.Rollback()
		return 0, err
	}

	//insert a record to liquidity_order and get the id
	state := 1
	if _, exist := config.SendCoins[coin1]; exist {
		state = 2
	}
	if _, exist := config.SendCoins[coin2]; exist {
		state = 2
	}
	var id int64
	if res, err := tx.Exec("insert into `liquidity_order`(`account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`) VALUES(?,?,?,?,?,?,-1,?,?)", account, coin1, coin2, amountOut1.String(), amountOut2.String(), amount.String(), state, time.Now().Unix()); err != nil {
		tx.Rollback()
		return 0, err
	} else {
		if id, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	//
	if _, exist := config.SendCoins[coin1]; exist {
		if _, err := tx.Exec("insert into `send_coin_pending`(`link_id`,`to`,`coin`,`amount`,`type`) VALUES(?,?,?,?,3)", id, account, coin1, amountOut1.String()); err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		amount1 := big.NewInt(amountOut1.Int64())
		if b, exist := balances[coin1]; exist {
			amount1.Add(amount1, b)
		}
		if _, err := tx.Exec("replace inito balance(account,coin,amount) values(?,?,?)", account, coin1, amount1.String()); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if _, exist := config.SendCoins[coin2]; exist {
		if _, err := tx.Exec("insert into `send_coin_pending`(`link_id`,`to`,`coin`,`amount`,`type`) VALUES(?,?,?,?,3)", id, account, coin2, amountOut2.String()); err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		amount2 := big.NewInt(amountOut2.Int64())
		if b, exist := balances[coin2]; exist {
			amount2.Add(amount2, b)
		}
		if _, err := tx.Exec("update balance set amount=? where account=? and coin=?", account, coin1, amount2.String()); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, err
}

//GetPendingLiquidityAddOrder
func GetPendingLiquidityAddOrder(account string) (PendingLiquidityAddOrder, error) {
	row := db.QueryRow("select `account`,`coin1`,`coin2`,`amount1`,`o_time` from liquidity_add_order_pending where `account`=?", account)

	o := PendingLiquidityAddOrder{}
	if err := row.Scan(&o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.OrderTime); err != nil {
		return o, err
	}
	return o, nil
}

//MovePendingLiquidityAddOrderToCancel
func MovePendingLiquidityAddOrderToCancel(account string) error {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return err
	}

	//lock row of table, query the pending liquidity add order
	row := tx.QueryRow("select `account`,`coin1`,`coin2`,`amount1`,`o_time` from liquidity_add_order_pending where `account`=? for update", account)
	o := PendingLiquidityAddOrder{}
	if err := row.Scan(&o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}

	//delete the record
	if _, err := tx.Exec("delete from liquidity_add_order_pending where account=?", account); err != nil {
		tx.Rollback()
		return err
	}

	//insert a record to liquidity_order
	a1, a2 := o.Amount1, "0"
	if o.Coin1 > o.Coin2 {
		o.Coin1, o.Coin2 = o.Coin2, o.Coin1
		a1, a2 = a2, a1
	}
	if _, err := tx.Exec("insert into `liquidity_order`(`account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`) VALUES(?,?,?,?,?,'0',1,4,?)", o.Account, o.Coin1, o.Coin2, a1, a2, o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//GetLiquidityOrder get liquidity order with id
func GetLiquidityOrder(id int64) (LiquidityOrder, error) {
	row := db.QueryRow("select `account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`,`e_time` from liquidity_order where id=?", id)

	o := LiquidityOrder{Id: id}
	if err := row.Scan(&o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.Amount2, &o.Lp, &o.Direction, &o.State, &o.OrderTime, &o.EndTime); err != nil {
		return o, err
	}
	return o, nil
}

//GetLiquidityOrders get liquidity orders by account
func GetLiquidityOrders(account string, count int) ([]LiquidityOrder, error) {
	rows, err := db.Query("select `id`,`account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`,`e_time` from liquidity_order where `account`=? order by id desc limit ?", account, count)
	if err != nil {
		return nil, err
	}

	os := make([]LiquidityOrder, 0)
	for rows.Next() {
		o := LiquidityOrder{}
		if err = rows.Scan(&o.Id, &o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.Amount2, &o.Lp, &o.Direction, &o.State, &o.OrderTime, &o.EndTime); err != nil {
			break
		}
		os = append(os, o)
	}

	return os, err
}
