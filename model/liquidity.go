package model

import "errors"

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

func InsertPendingLiquidityOrder(account, coin1, coin2, amount string, direction int) error {
	if _, err := db.Exec("INSERT INTO `liquidity_order_pending`(`account`,`coin1`,`coin2`,`amount`,`direction`) VALUES(?,?,?,?,?)", account, coin1, coin2, amount, direction); err != nil {
		return err
	}
	return nil
}

func GetPendingLiquidityOrder(account string) (LiquidityOrder, error) {
	row := db.QueryRow("select `account`,`coin1`,`coin2`,`amount`,`direction`,`o_time` from liquidity_order_pending where `account`=?", account)

	o := LiquidityOrder{}
	if err := row.Scan(&o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.Direction, &o.OrderTime); err != nil {
		return o, err
	}
	if o.Direction == -1 {
		o.Lp = o.Amount1
		o.Amount1 = ""
	}
	return o, nil
}

func MovePendingLiquidityOrderToCancel(account string) error {
	o, err := GetPendingLiquidityOrder(account)
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
	res, err := tx.Exec("delete from liquidity_order_pending where account=?", account)
	if err != nil {
		tx.Rollback()
		return err
	}
	if c, err := res.RowsAffected(); err != nil {
		tx.Rollback()
		return err
	} else if c == 0 {
		tx.Rollback()
		return errors.New("have no liquidity_order_pending")
	}
	if _, err := tx.Exec("INSERT INTO `liquidity_order`(`account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,state`,`o_time`) VALUES(?,?,?,?,?,?,?,?,?)", o.Account, o.Coin1, o.Coin2, o.Amount1, "", o.Lp, o.Direction, 4, o.OrderTime); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetLiquidityOrder(id int64) (LiquidityOrder, error) {
	row := db.QueryRow("select `account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`,`e_time` from liquidity_order where id=?", id)

	o := LiquidityOrder{Id: id}
	if err := row.Scan(&o.Account, &o.Coin1, &o.Coin2, &o.Amount1, &o.Amount2, &o.Lp, &o.Direction, &o.State, &o.OrderTime, &o.EndTime); err != nil {
		return o, err
	}
	return o, nil
}

func GetLiquidityOrders(account string, count int) ([]LiquidityOrder, error) {
	rows, err := db.Query("select `id`,`account`,`coin1`,`coin2`,`amount1`,`amount2`,`lp`,`direction`,`state`,`o_time`,`e_time` from liquidity_order where `account`=? order by id desc limit ?", account, count)

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
