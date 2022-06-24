package model

import (
	"errors"
	"math/big"
)

type PendingCollectOrder struct {
	Account   string `json:"account"`
	Address   string `json:"address"`
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	OrderTime int    `json:"o_time"`
}

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

func InsertPendingCollectOrder(account, address, coin, amount string) error {
	if _, err := db.Exec("insert into `collect_order_pending`(`account`,`address`,`coin`,`amount`) VALUES(?,?,?,?)", account, address, coin, amount); err != nil {
		return err
	}
	return nil
}

func RetrieveCoin(account, coin, amount, to string) error {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return err
	}

	//get and judge the balance
	row := tx.QueryRow("select `amount` from balance where `account`=? and `coin`=? for update", account, coin)
	var balance string
	if err := row.Scan(&balance); err != nil {
		tx.Rollback()
		return err
	}
	balanceUser, b1 := new(big.Int).SetString(balance, 10)
	a, b2 := new(big.Int).SetString(amount, 10)
	if !b1 || !b2 || balanceUser.Cmp(a) < 0 {
		tx.Rollback()
		return errors.New("balance is not enough : " + balance)
	}

	if _, err := tx.Exec("update balance set amount=? where account=? and coin=?", balanceUser.Sub(balanceUser, a).String(), account, coin); err != nil {
		tx.Rollback()
		return err
	}

	var id int64
	if res, err := db.Exec("insert into `coin_order`(`account`,`address`,`coin`,`amount`,`direction`,`state`) VALUES(?,?,?,?,-1,2)", account, to, coin, amount); err != nil {
		tx.Rollback()
		return err
	} else {
		if id, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, err := tx.Exec("insert into `send_coin_pending`(`link_id`,`to`,`coin`,`amount`,`type`) VALUES(?,?,?,?,2)", id, to, coin, amount); err != nil {
		return err
	}

	return tx.Commit()
}

func GetPendingCollectOrder(address string) (PendingCollectOrder, error) {
	row := db.QueryRow("select `account`,`address`,`coin`,`amount`,`o_time` from collect_order_pending where `account`=? or `address`=?", address, address)

	o := PendingCollectOrder{}
	if err := row.Scan(&o.Account, &o.Address, &o.Coin, &o.Amount, &o.OrderTime); err != nil {
		return o, err
	}
	return o, nil
}

func MovePendingCollectOrderToCancel(address string) error {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return err
	}

	row := tx.QueryRow("select `account`,`address`,`coin`,`amount`,`o_time` from collect_order_pending where `account`=? or `address`=? for update", address, address)
	o := PendingCollectOrder{}
	if err := row.Scan(&o.Account, &o.Address, &o.Coin, &o.Amount, &o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("delete from collect_order_pending where address=? or account=?", address, address); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("insert into `coin_order`(`account`,`address`,`coin`,`amount`,`direction`,`state`,`o_time`) VALUES(?,?,?,?,1,4,?)", o.Account, o.Address, o.Coin, o.Amount, 1, 4, o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetCoinOrder(id int64) (CoinOrder, error) {
	row := db.QueryRow("select `id`,`account`,`address`,`coin`,`amount`,`direction`,`state`,`o_time`,`e_time` from coin_order where id=?", id)

	o := CoinOrder{}
	if err := row.Scan(&o.Id, &o.Account, &o.Address, &o.Coin, &o.Amount, &o.Direction, &o.State, &o.OrderTime, &o.EndTime); err != nil {
		return o, err
	}
	return o, nil
}

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
