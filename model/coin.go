package model

import "errors"

type CoinOrder struct {
	Id        int64  `json:"id"`
	Account   string `json:"account"`
	From      string `json:"from"`
	Address   string `json:"address"`
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	State     int    `json:"state"`
	Direction int    `json:"direction"`
	OrderTime int    `json:"o_time"`
	EndTime   int    `json:"e_time"`
}

func InsertPendingCoinOrder(account, address, coin, amount string, direction int) error {
	if _, err := db.Exec("INSERT INTO `coin_order_pending`(`account`,`address`,`coin`,`amount`,`direction`) VALUES(?,?,?,?,?)", account, address, coin, amount, direction); err != nil {
		return err
	}
	return nil
}

func GetPendingCoinOrder(address string) (CoinOrder, error) {
	row := db.QueryRow("select `account`,`address`,`coin`,`amount`,`direction`,`o_time` from coin_order_pending where `account`=? or `address`=?", address, address)

	o := CoinOrder{}
	if err := row.Scan(&o.Account, &o.Address, &o.Coin, &o.Amount, &o.Direction, &o.OrderTime); err != nil {
		return o, err
	}
	return o, nil
}

func MovePendingCoinOrderToCancel(address string) error {
	o, err := GetPendingCoinOrder(address)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return err
	}
	res, err := tx.Exec("delete from coin_order_pending where (direction=1 and address=?) or (direction=-1 and account=?)", address, address)
	if err != nil {
		tx.Rollback()
		return err
	}
	if c, err := res.RowsAffected(); err != nil {
		tx.Rollback()
		return err
	} else if c == 0 {
		tx.Rollback()
		return errors.New("have no coin_order_pending")
	}
	if _, err := tx.Exec("INSERT INTO `coin_order`(`account`,`address`,`coin`,`amount`,`direction`,`stata`,`o_time`) VALUES(?,?,?,?,?,?,?)", o.Account, o.Address, o.Coin, o.Amount, o.Direction, 4, o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetCoinOrder(id int64) (CoinOrder, error) {
	row := db.QueryRow("select `id`,`account`,`address`,`coin`,`amount`,`direction`,`state`,`o_time`,`e_time` from collect_order where id=?", id)

	o := CoinOrder{}
	if err := row.Scan(&o.Id, &o.Account, &o.Address, &o.Coin, &o.Amount, &o.Direction, &o.State, &o.OrderTime, &o.EndTime); err != nil {
		return o, err
	}
	return o, nil
}

func GetCoinOrders(address string, count int) ([]CoinOrder, error) {
	rows, err := db.Query("select `id`,`account`,`address`,`coin`,`amount`,`direction`,`state`,`o_time`,`e_time` from collect_order where `address`=? or account=? order by id desc limit ?", address, address, count)

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
