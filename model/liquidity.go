package model

import (
	"database/sql"
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
	Amount2   string `json:"amount2"`
	OrderTime int    `json:"o_time"`
}

type LiquidityAddOrder struct {
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

func InsertPendingLiquidityAddOrder(account, coin1, coin2, amount1, amount2 string) error {
	if _, err := db.Exec("insert into `liquidity_add_order_pending`(`account`,`coin1`,`coin2`,`amount1`,`amount2`) VALUES(?,?,?,?,?)", account, coin1, coin2, amount1, amount2); err != nil {
		return err
	}
	return nil
}

func RemoveLiquidity(account, coin1, coin2, amount string) error {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return err
	}
	lpCoin := "LP-" + coin1 + "-" + coin2
	removeLiquidity, b := new(big.Int).SetString(amount, 10)
	if !b {
		tx.Rollback()
		return fmt.Errorf("removeLiquidity convert To big.Int error. %s", amount)
	}

	row := tx.QueryRow("select amount from balance where account=? and coin=? for update", account, lpCoin)
	var str string
	if err = row.Scan(&str); err != nil {
		tx.Rollback()
		return err
	}
	liquidity, b := new(big.Int).SetString(str, 10)
	if !b || liquidity.Cmp(removeLiquidity) < 0 {
		tx.Rollback()
		return fmt.Errorf("liquidity convert To big.Int error or balance is not enough. %s : %s", str, amount)
	}

	row = tx.QueryRow("select reserve1,reserve2,total_supply from swap_pair where coin1=? and coin2=? for update", coin1, coin2)
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

	amountOut1 := new(big.Int).Div(new(big.Int).Mul(removeLiquidity, reserve1), totalSupply)
	amountOut2 := new(big.Int).Div(new(big.Int).Mul(removeLiquidity, reserve2), totalSupply)
	if amountOut1.Cmp(big.NewInt(0)) <= 0 || amountOut2.Cmp(big.NewInt(0)) <= 0 {
		tx.Rollback()
		return fmt.Errorf("insufficient liquidity burned")
	}

	if _, err := tx.Exec("update swap_pair set amount1=?,amount2=?,total_supply=? where coin1=? and coin2=?", reserve1.Sub(reserve1, amountOut1).String(), reserve2.Sub(reserve2, amountOut2).String(), totalSupply.Sub(totalSupply, removeLiquidity).String(), coin1, coin2); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("update balance set amount=? where account=? and coin=?", liquidity.Sub(liquidity, removeLiquidity).String(), account, amount); err != nil {
		tx.Rollback()
		return err
	}

	state := 1
	if _, exist := config.SendCoins[coin1]; exist {
		state = 2
	}
	if _, exist := config.SendCoins[coin2]; exist {
		state = 2
	}
	var id int64
	if res, err := tx.Exec("insert into `liquidity_order`(`account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`) VALUES(?,?,?,?,?,?,-1,2,?)", account, coin1, coin2, amountOut1.String(), amountOut2.String(), amount, state, time.Now().Unix()); err != nil {
		tx.Rollback()
		return err
	} else {
		if id, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, exist := config.SendCoins[coin1]; exist {
		if _, err := tx.Exec("insert into `send_coin_pending`(`link_id`,`to`,`coin`,`amount`,`type`) VALUES(?,?,?,?,3)", id, account, coin2, amountOut2.String()); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		row = tx.QueryRow("select amount from balance where account=? and coin=?", account, coin1)
		if err = row.Scan(&str); err != nil {
			if err != sql.ErrNoRows {
				tx.Rollback()
				return err
			}
			str = "0"
		}
		amount1, b := new(big.Int).SetString(str, 10)
		if !b {
			tx.Rollback()
			return fmt.Errorf("amount1 convert To big.Int error. %s", str)
		}

		if _, err := tx.Exec("update balance set amount=? where account=? and coin=?", amount1.Add(amount1, amountOut1), account, coin1); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, exist := config.SendCoins[coin2]; exist {
		if _, err := tx.Exec("insert into `send_coin_pending`(`link_id`,`to`,`coin`,`amount`,`type`) VALUES(?,?,?,?,3)", id, account, coin2, amountOut2.String()); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		row = tx.QueryRow("select amount from balance where account=? and coin=?", account, coin2)
		if err = row.Scan(&str); err != nil {
			if err != sql.ErrNoRows {
				tx.Rollback()
				return err
			}
			str = "0"
		}
		amount2, b := new(big.Int).SetString(str, 10)
		if !b {
			tx.Rollback()
			return fmt.Errorf("amount2 convert To big.Int error. %s", str)
		}

		if _, err := tx.Exec("update balance set amount=? where account=? and coin=?", amount2.Add(amount2, amountOut1), account, coin1); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetPendingLiquidityAddOrder(account string) (PendingLiquidityAddOrder, error) {
	row := db.QueryRow("select `account`,`coin1`,`coin2`,`amount1`,`amount2`,`o_time` from liquidity_order_add_pending where `account`=?", account)

	o := PendingLiquidityAddOrder{}
	if err := row.Scan(&o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.Amount2, &o.OrderTime); err != nil {
		return o, err
	}
	return o, nil
}

func MovePendingLiquidityAddOrderToCancel(account string) error {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return err
	}

	row := tx.QueryRow("select `account`,`coin1`,`coin2`,`amount1`,`amount2`,`o_time` from liquidity_order_pending where `account`=? for update", account)
	o := PendingLiquidityAddOrder{}
	if err := row.Scan(&o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.Amount2, &o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("delete from liquidity_order_pending where account=?", account); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec("INSERT INTO `liquidity_order`(`account`,`coin1`,`coin2`,`reserve1`,`reserver2`,`total_supply`,`direction`,`state`,`o_time`) VALUES(?,?,?,?,?,?,?,?,?)", o.Account, o.Coin1, o.Coin2, o.Amount1, o.Amount2, "0", 1, 4, o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetLiquidityOrder(id int64) (LiquidityAddOrder, error) {
	row := db.QueryRow("select `account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`,`e_time` from liquidity_order where id=?", id)

	o := LiquidityAddOrder{Id: id}
	if err := row.Scan(&o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.Amount2, &o.Lp, &o.Direction, &o.State, &o.OrderTime, &o.EndTime); err != nil {
		return o, err
	}
	return o, nil
}

func GetLiquidityOrders(account string, count int) ([]LiquidityAddOrder, error) {
	rows, err := db.Query("select `id`,`account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`,`e_time` from liquidity_order where `account`=? order by id desc limit ?", account, count)
	if err != nil {
		return nil, err
	}

	os := make([]LiquidityAddOrder, 0)
	for rows.Next() {
		o := LiquidityAddOrder{}
		if err = rows.Scan(&o.Id, &o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.Amount2, &o.Lp, &o.Direction, &o.State, &o.OrderTime, &o.EndTime); err != nil {
			break
		}
		os = append(os, o)
	}

	return os, err
}
