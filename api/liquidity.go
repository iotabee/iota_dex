package api

import (
	"iota_dex/config"
	"iota_dex/gl"
	"iota_dex/model"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func LiquidityAddOrder(c *gin.Context) {
	account := c.GetString("account")
	coin1 := strings.ToUpper(c.Query("coin1"))
	coin2 := strings.ToUpper(c.Query("coin2"))
	amount1 := c.Query("amount1")

	_, _, _, err := model.GetPrice(coin1, coin2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "have no pair",
		})
		gl.OutLogger.Error("Get price when add liquidity order error. %s, %s, %v", coin1, coin2, err)
		return
	}

	a, b1 := new(big.Int).SetString(amount1, 10)
	if !b1 {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "params error.",
		})
		gl.OutLogger.Error("Add liquidity order params error. %s", amount1)
		return
	}

	if _, exist := config.SendCoins[coin1]; exist {
		err = model.InsertPendingLiquidityAddOrder(account, coin1, coin2, amount1)
	} else {
		err = model.AddLiquidity(account, coin1, coin2, a)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "maybe you have a liquidity add order is pending.",
		})
		gl.OutLogger.Error("add liquidity error. %s, %s, %s, %s, %v", account, coin1, coin2, amount1, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func LiquidityRemoveOrder(c *gin.Context) {
	account := c.GetString("account")
	coin1 := strings.ToUpper(c.Query("coin1"))
	coin2 := strings.ToUpper(c.Query("coin2"))
	lp, b := new(big.Int).SetString(c.Query("lp"), 10)

	if coin1 > coin2 {
		coin1, coin2 = coin2, coin1
	}
	if _, _, _, err := model.GetPrice(coin1, coin2); err != nil || !b {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "have no pair",
		})
		gl.OutLogger.Error("remove liquidity request params error. %s, %s, %v, %v", coin1, coin2, c.Query("lp"), err)
		return
	}

	if err := model.RemoveLiquidity(account, coin1, coin2, lp); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "maybe balance is not enough",
		})
		gl.OutLogger.Error("Remove liquidity in db error. %s, %s, %s, %s, %v", account, coin1, coin2, lp, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func CancelPendingLiquidityAddOrder(c *gin.Context) {
	account := c.GetString("account")
	err := model.MovePendingLiquidityAddOrderToCancel(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending liquidity add order",
		})
		gl.OutLogger.Error("cancel liquidity_order_add_pending error. %s, %v", account, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func GetPendingLiquidityAddOrder(c *gin.Context) {
	account := c.GetString("account")
	o, err := model.GetPendingLiquidityAddOrder(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending liquidity add order",
		})
		gl.OutLogger.Error("get liquidity_order_add_pending error. %s, %v", account, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}

func GetLiquidityOrders(c *gin.Context) {
	account := c.GetString("account")
	count, _ := strconv.Atoi(c.DefaultQuery("count", "5"))
	if count == 0 {
		count = 5
	}
	o, err := model.GetLiquidityOrders(account, count)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no liquidity orders",
		})
		gl.OutLogger.Error("get liquidity_order error. %s, %s, %v", account, c.Query("count"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}

func GetLiquidityOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "order id error",
		})
		gl.OutLogger.Error("liquidity order id error. %d, %v", id, err)
		return
	}

	o, err := model.GetLiquidityOrder(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no this order",
		})
		gl.OutLogger.Error("get liquidity_order error. %d, %v", id, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}
