package api

import (
	"iota_dex/gl"
	"iota_dex/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func LiquidityAddOrder(c *gin.Context) {
	account := c.GetString("account")
	coin := strings.ToUpper(c.Query("coin"))
	amount := c.Query("amount")
	coin1 := strings.ToUpper(c.Query("coin1"))

	c1, c2 := coin, coin1
	if c1 > c2 {
		c1, c2 = c2, c1
	}
	if _, _, _, err := model.GetPrice(c1, c2); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "have no pair",
		})
		gl.OutLogger.Error("Add liquidity order error. %s, %s, %v", c1, c2, err)
		return
	}
	if len(amount) == 0 || len(account) < 0 {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "params error.",
		})
		gl.OutLogger.Error("Add liquidity order params error. %s : %s", amount, account)
		return
	}

	if err := model.InsertPendingLiquidityOrder(account, coin, coin1, amount, 1); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "maybe you have a liquidity order is pending.",
		})
		gl.OutLogger.Error("Insert into db error(pending_liquidity_order). 1, %s, %s, %s, %s, %v", account, coin, coin1, amount, err)
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
	lp := c.Query("lp")

	if coin1 > coin2 {
		coin1, coin2 = coin2, coin1
	}
	if _, _, _, err := model.GetPrice(coin1, coin2); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "have no pair",
		})
		gl.OutLogger.Error("Remove liquidity order error. %s, %s, %v", coin1, coin2, err)
		return
	}
	if len(lp) == 0 || len(account) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "params error.",
		})
		gl.OutLogger.Error("Remove liquidity order params error. %s : %s", lp, account)
		return
	}

	if err := model.InsertPendingLiquidityOrder(account, coin1, coin2, lp, -1); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "maybe you have a liquidity order is pending.",
		})
		gl.OutLogger.Error("Insert into db error(liquidity_order_pending). -1, %s, %s, %s, %s, %v", account, coin1, coin2, lp, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func CancelPendingLiquidityOrder(c *gin.Context) {
	account := c.GetString("account")
	err := model.MovePendingLiquidityOrderToCancel(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending liquidity order",
		})
		gl.OutLogger.Error("cancel liquidity_order_pending error. %s, %v", account, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func GetPendingLiquidityOrder(c *gin.Context) {
	account := c.GetString("account")
	o, err := model.GetPendingLiquidityOrder(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending liquidity order",
		})
		gl.OutLogger.Error("get liquidity_order_pending error. %s, %v", account, err)
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
