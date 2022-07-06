package model

type Coin struct {
	Contract string `json:"contract"`
	Wallet   string `json:"wallet"`
	Decimal  int    `json:"decimal"`
	Amount   string `json:"amount"`
}

type SwapPair struct {
	Id          int
	Coin1       string  `json:"coin1"`
	Coin2       string  `json:"coin2"`
	Reserve1    string  `json:"reserve1"`
	Reserve2    string  `json:"reserve2"`
	TotalSupply string  `json:"total_supply"`
	FeeRate     float32 `json:"fee_rate"`
	FeeScale    float32 `json:"fee_scale"`
}

func GetPairs() ([]SwapPair, error) {
	rows, err := db.Query("select id,coin1,coin2,reserve1,reserve2,total_supply,fee_rate,fee_scale from swap_pair")
	if err != nil {
		return nil, err
	}

	pairs := make([]SwapPair, 0)
	for rows.Next() {
		sp := SwapPair{}
		rows.Scan(&sp.Id, &sp.Coin1, &sp.Coin2, &sp.Reserve1, &sp.Reserve2, &sp.TotalSupply, &sp.FeeRate, &sp.FeeScale)
		pairs = append(pairs, sp)
	}

	return pairs, nil
}

func GetPair(coin1, coin2 string) (SwapPair, error) {
	row := db.QueryRow("select coin1,coin2,reserve1,reserve2,total_supply,fee_rate,fee_scale from swap_pair where coin1=? and coin2=?", coin1, coin2)
	sp := SwapPair{}
	err := row.Scan(&sp.Coin1, &sp.Coin2, &sp.Reserve1, &sp.Reserve2, &sp.TotalSupply, &sp.FeeRate, &sp.FeeScale)
	return sp, err
}

func GetCoin(symbol string) (Coin, error) {
	row := db.QueryRow("select `contract`,`deci`,`wallet` from coin where symbol=?", symbol)
	c := Coin{}
	err := row.Scan(&c.Contract, &c.Decimal, &c.Wallet)
	return c, err
}
