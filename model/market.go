package model

type Coin struct {
	Contract string `json:"contract"`
	Wallet   string `json:"wallet"`
	Decimal  int    `json:"decimal"`
	Amount   string `json:"amount"`
}

type SwapPair struct {
	Id      int             `json:"id"`
	Coins   map[string]Coin `json:"coins"`
	Lp      string          `json:"lp"`
	FeeRate float32         `json:"fee_rate"`
}

func GetPairs() ([]SwapPair, error) {
	rows, err := db.Query("select symbol,contract,deci,wallet from coin")
	if err != nil {
		return nil, err
	}
	coins := make(map[string]Coin)
	for rows.Next() {
		c := Coin{}
		symbol := ""
		rows.Scan(&symbol, &c.Contract, &c.Decimal, &c.Wallet)
		coins[symbol] = c
	}

	rows, err = db.Query("select id,coin1,coin2,amount1,amount2,lp,fee_rate from swap_pair")
	if err != nil {
		return nil, err
	}

	pairs := make([]SwapPair, 0)
	for rows.Next() {
		sp := SwapPair{}
		sp.Coins = make(map[string]Coin)
		var coin1, coin2, amount1, amount2 string
		rows.Scan(&sp.Id, &coin1, &coin2, &amount1, &amount2, &sp.Lp, &sp.FeeRate)
		if c, exist := coins[coin1]; exist {
			c.Amount = amount1
			sp.Coins[coin1] = c
		}
		if c, exist := coins[coin2]; exist {
			c.Amount = amount2
			sp.Coins[coin2] = c
		}
		pairs = append(pairs, sp)
	}

	return pairs, nil
}

func GetPrice(coin1, coin2 string) (string, string, float32, error) {
	row := db.QueryRow("select amount1,amount2,fee_rate from swap_pair where coin1=? and coin2=?", coin1, coin2)
	var a1, a2 string
	var fr float32
	if err := row.Scan(&a1, &a2, &fr); err != nil {
		return "", "", 0, err
	}
	return a1, a2, fr, nil
}

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
