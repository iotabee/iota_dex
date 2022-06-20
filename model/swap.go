package model

import (
	"errors"
)

type SwapOrder struct {
	Id         int64  `json:"id"`
	FromAddr   string `json:"from_address"`
	FromCoin   string `json:"from_coin"`
	FromAmount string `json:"from_amount"`
	ToAddr     string `json:"to_address"`
	ToCoin     string `json:"to_coin"`
	ToAmount   string `json:"to_amount"`
	State      int    `json:"state"`
	OrderTime  int    `json:"o_time"`
	EndTime    int    `json:"e_time"`
}

func InsertPendingSwapOrder(fromAddr, fromCoin, fromAmount, toAddr, toCoin, minAmount string) error {
	if _, err := db.Exec("INSERT INTO `swap_order_pending`(`from_address`,`from_coin`,`from_amount`,`to_address`,`to_coin`,`min_amount`) VALUES(?,?,?,?,?,?)", fromAddr, fromCoin, fromAmount, toAddr, toCoin, minAmount); err != nil {
		return err
	}
	return nil
}

func GetPendingSwapOrder(account string) (SwapOrder, error) {
	row := db.QueryRow("select `from_address`,`from_coin`,`from_amount`,`to_address`,`to_coin`,`min_amount`,`o_time` from swap_order_pending where `from_address`=?", account)

	o := SwapOrder{State: 0}
	if err := row.Scan(&o.FromAddr, &o.FromCoin, &o.FromAmount, &o.ToAddr, &o.ToCoin, &o.ToAmount, &o.OrderTime); err != nil {
		return o, err
	}
	return o, nil
}

func MovePendingSwapOrderToCancel(account string) error {
	o, err := GetPendingSwapOrder(account)
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
	res, err := tx.Exec("delete from swap_order_pending where from_address=?", account)
	if err != nil {
		tx.Rollback()
		return err
	}
	if c, err := res.RowsAffected(); err != nil {
		tx.Rollback()
		return err
	} else if c == 0 {
		tx.Rollback()
		return errors.New("have no swap_order_pending")
	}
	if _, err := tx.Exec("INSERT INTO `swap_order`(`from_address`,`from_coin`,`from_amount`,`to_address`,`to_coin`,`to_amount`,`state`,`o_time`) VALUES(?,?,?,?,?,?,?,?)", o.FromAddr, o.FromCoin, o.FromAmount, o.ToAddr, o.ToCoin, o.ToAmount, 4, o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetSwapOrder(id int64) (SwapOrder, error) {
	row := db.QueryRow("select `id`,`from_address`,`from_coin`,`from_amount`,`to_address`,`to_coin`,`to_amount`,`state`,`o_time`,`e_time` from swap_order where id=?", id)

	o := SwapOrder{}
	if err := row.Scan(&o.Id, &o.FromAddr, &o.FromCoin, &o.FromAmount, &o.ToAddr, &o.ToCoin, &o.ToAmount, &o.State, &o.OrderTime, &o.EndTime); err != nil {
		return o, err
	}
	return o, nil
}

func GetSwapOrders(account string, count int) ([]SwapOrder, error) {
	rows, err := db.Query("select `id`,`from_address`,`from_coin`,`from_amount`,`to_address`,`to_coin`,`to_amount`,`state`,`o_time`,`e_time` from swap_order where `from_address`=? order by id desc limit ?", account, count)
	if err != nil {
		return nil, err
	}

	os := make([]SwapOrder, 0)
	for rows.Next() {
		o := SwapOrder{}
		if err = rows.Scan(&o.Id, &o.FromAddr, &o.FromCoin, &o.FromAmount, &o.ToAddr, &o.ToCoin, &o.ToAmount, &o.State, &o.OrderTime, &o.EndTime); err != nil {
			break
		}
		os = append(os, o)
	}

	return os, err
}
