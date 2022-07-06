package model

import (
	"errors"
	"math/big"
	"time"
)

//GetBalance get account's balance
func GetBalance(account string) (map[string]string, error) {
	rows, err := db.Query("select coin,amount from balance where account=?", account)
	if err != nil {
		return nil, err
	}
	b := make(map[string]string)
	for rows.Next() {
		var coin, amount string
		err = rows.Scan(&coin, &amount)
		if err != nil {
			return nil, err
		}
		b[coin] = amount
	}
	return b, nil
}

//PendingCollectOrder, collect coin order which is pending
type PendingCollectOrder struct {
	Account   string `json:"account"`
	From      string `json:"from"`
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	OrderTime int    `json:"o_time"`
}

//CoinOrder, collect coin or retrieve coin order
type CoinOrder struct {
	Id        int64  `json:"id"`
	Account   string `json:"account"`
	Address   string `json:"address"`
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	Direction int    `json:"direction"`
	State     int    `json:"state"`
	OrderTime int    `json:"o_time"`
	EndTime   int    `json:"e_time"`
}

//InsertPendingCollectOrder, insert a collect_order_pending record to db.
//@account 	: address that the coin to collect
//@from 	: address which the coin from
//@coin		: the coin type as string of upper
//@amount	: amount that the coin to collect
//@return	: nil if successful, or error if failed
func InsertPendingCollectOrder(account, from, coin, amount string) error {
	if _, err := db.Exec("insert into `collect_order_pending`(`account`,`from`,`coin`,`amount`) VALUES(?,?,?,?)", account, from, coin, amount); err != nil {
		return err
	}
	return nil
}

//RetrieveCoin, retrieve coin from the account. Use transaction of db to deal the balance of account.
//@account	: address that the coin to retrieve
//@coin		: the coin type as string of upper
//@amount	: amount that the coin to retrieve
//@to		: address which the coin to send
//@return	: nil if successful, or error if failed
func RetrieveCoin(account, coin, to string, amount *big.Int) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return 0, err
	}

	//Get balance from db, and lock the 'balance' table
	row := tx.QueryRow("select `amount` from balance where `account`=? and `coin`=? for update", account, coin)
	var str string
	if err := row.Scan(&str); err != nil {
		tx.Rollback()
		return 0, err
	}
	balance, _ := new(big.Int).SetString(str, 10)
	if balance.Cmp(amount) < 0 {
		tx.Rollback()
		return 0, errors.New("balance is not enough. " + str + " : " + amount.String())
	}

	//update the balance
	if _, err := tx.Exec("update balance set amount=? where account=? and coin=?", balance.Sub(balance, amount).String(), account, coin); err != nil {
		tx.Rollback()
		return 0, err
	}

	//insert a record to coin_order of db and get the id
	var id int64
	if res, err := db.Exec("insert into `coin_order`(`account`,`address`,`coin`,`amount`,`direction`,`state`,`o_time`) VALUES(?,?,?,?,-1,2,?)", account, to, coin, amount.String(), time.Now().Unix()); err != nil {
		tx.Rollback()
		return 0, err
	} else {
		if id, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	//insert a record to send_coin_pending of db, set link_id=id
	if _, err := tx.Exec("insert into `send_coin_pending`(`link_id`,`to`,`coin`,`amount`,`type`) VALUES(?,?,?,?,2)", id, to, coin, amount.String()); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

//GetPendingCollectOrder, get the pending collect order from db
// @address	: address of account or from
//@return	: PendingCollectOrder as struct
func GetPendingCollectOrder(address string) (PendingCollectOrder, error) {
	row := db.QueryRow("select `account`,`from`,`coin`,`amount`,`o_time` from collect_order_pending where `account`=? or `from`=?", address, address)
	o := PendingCollectOrder{}
	if err := row.Scan(&o.Account, &o.From, &o.Coin, &o.Amount, &o.OrderTime); err != nil {
		return o, err
	}
	return o, nil
}

//MovePendingCollectOrderToCancel, Cancel the pending collect order, move the record of collect_order_pending to coin_order table.
func MovePendingCollectOrderToCancel(address string) error {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return err
	}

	//query the pending collect order and lock the table
	row := tx.QueryRow("select `account`,`from`,`coin`,`amount`,`o_time` from collect_order_pending where `account`=? or `from`=? for update", address, address)
	o := PendingCollectOrder{}
	if err := row.Scan(&o.Account, &o.From, &o.Coin, &o.Amount, &o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}

	//delete the record of collect_order_pending.
	if _, err := tx.Exec("delete from collect_order_pending where `account`=? or `from`=?", address, address); err != nil {
		tx.Rollback()
		return err
	}

	//add a record of coin_order
	if _, err := tx.Exec("insert into `coin_order`(`account`,`address`,`coin`,`amount`,`direction`,`state`,`o_time`) VALUES(?,?,?,?,1,4,?)", o.Account, o.From, o.Coin, o.Amount, o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

//GetCoinOrders, the recent coin orders by the account or address
func GetCoinOrders(address string, count int) ([]CoinOrder, error) {
	rows, err := db.Query("select `id`,`account`,`address`,`coin`,`amount`,`direction`,`state`,`o_time`,`e_time` from coin_order where `address`=? or account=? order by id desc limit ?", address, address, count)
	if err != nil {
		return nil, err
	}
	os := make([]CoinOrder, 0)
	for rows.Next() {
		o := CoinOrder{}
		if err = rows.Scan(&o.Id, &o.Account, &o.Address, &o.Coin, &o.Amount, &o.Direction, &o.State, &o.OrderTime, &o.EndTime); err != nil {
			break
		}
		os = append(os, o)
	}

	return os, err
}

//GetCoinOrder, the coin order with id.
func GetCoinOrder(id int64) (CoinOrder, error) {
	row := db.QueryRow("select `id`,`account`,`address`,`coin`,`amount`,`direction`,`state`,`o_time`,`e_time` from coin_order where id=?", id)
	o := CoinOrder{}
	if err := row.Scan(&o.Id, &o.Account, &o.Address, &o.Coin, &o.Amount, &o.Direction, &o.State, &o.OrderTime, &o.EndTime); err != nil {
		return o, err
	}
	return o, nil
}
