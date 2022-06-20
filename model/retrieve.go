package model

import "errors"

type RetrieveOrder struct {
	Id        int64  `json:"id"`
	Account   string `json:"account"`
	To        string `json:"to"`
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	State     int    `json:"state"`
	OrderTime int    `json:"o_time"`
	EndTime   int    `json:"e_time"`
}

func InsertPendingRetrieveOrder(account, to, coin, amount string) error {
	if _, err := db.Exec("INSERT INTO `pending_retrieve_order`(`account`,`to`,`coin`,`amount`) VALUES(?,?,?,?)", account, to, coin, amount); err != nil {
		return err
	}
	return nil
}

func GetPendingRetrieveOrder(account string) (RetrieveOrder, error) {
	row := db.QueryRow("select `account`,`to`,`coin`,`amount`,`o_time` from pending_retrieve_order where `account`=?", account)

	o := RetrieveOrder{}
	if err := row.Scan(&o.Account, &o.To, &o.Coin, &o.Amount, &o.OrderTime); err != nil {
		return o, err
	}
	return o, nil
}

func MovePendingRetrieveOrderToCancel(account string) error {
	o, err := GetPendingRetrieveOrder(account)
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
	res, err := tx.Exec("delete from pending_retrieve_order where account=?", account)
	if err != nil {
		tx.Rollback()
		return err
	}
	if c, err := res.RowsAffected(); err != nil {
		tx.Rollback()
		return err
	} else if c == 0 {
		tx.Rollback()
		return errors.New("have no pending_retrieve_order")
	}
	if _, err := tx.Exec("INSERT INTO `retrieve_order`(`account`,`to`,`coin`,`amount`,`stata`,`o_time`) VALUES(?,?,?,?,?,?)", o.Account, o.To, o.Coin, o.Amount, 4, o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetRetrieveOrder(id int64) (RetrieveOrder, error) {
	row := db.QueryRow("select `id`,`account`,`to`,`coin`,`amount`,`state`,`o_time`,`e_time` from retrieve_order where id=?", id)

	o := RetrieveOrder{}
	if err := row.Scan(&o.Id, &o.Account, &o.To, &o.Coin, &o.Amount, &o.State, &o.OrderTime, &o.EndTime); err != nil {
		return o, err
	}
	return o, nil
}

func GetRetrieveOrders(account string, count int) ([]RetrieveOrder, error) {
	rows, err := db.Query("select `id`,`account`,`to`,`coin`,`amount`,`state`,`o_time`,`e_time` from retrieve_order where `account`=? order by id desc limit ?", account, count)
	if err != nil {
		return nil, err
	}

	os := make([]RetrieveOrder, 0)
	for rows.Next() {
		o := RetrieveOrder{}
		if err = rows.Scan(&o.Id, &o.Account, &o.To, &o.Coin, &o.Amount, &o.State, &o.OrderTime, &o.EndTime); err != nil {
			break
		}
		os = append(os, o)
	}

	return os, err
}
