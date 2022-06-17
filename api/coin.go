package api

import (
	"iota_dex/gl"
	"iota_dex/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CollectCoinOrder(c *gin.Context) {
	coin := c.Query("coin")
	amount := c.Query("amount")
	account := c.Query("account")

	coin = strings.ToUpper(coin)
	if len(coin) == 0 || len(account) == 0 || len(amount) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "params error",
		})
		gl.OutLogger.Error("param error when colect coin. %s, %s, %s", coin, account, amount)
		return
	}

	if err := model.InsertPendingCoinOrder(account, c.GetString("account"), coin, amount, 1); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "maybe you have a coin order is pending.",
		})
		gl.OutLogger.Error("Insert into db error(coin_order_pending). %s, %s, %s, %s, %v", account, c.GetString("account"), coin, amount, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func RetrieveCoinOrder(c *gin.Context) {
	to := c.Query("to")
	coin := c.Query("coin")
	amount := c.Query("amount")

	coin = strings.ToUpper(coin)
	if len(coin) == 0 || len(to) == 0 || len(amount) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "params error",
		})
		gl.OutLogger.Error("param error when retrieve coin. %s, %s, %s", to, coin, amount)
		return
	}

	if err := model.InsertPendingCoinOrder(c.GetString("account"), to, coin, amount, -1); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "maybe you have a coin order is pending.",
		})
		gl.OutLogger.Error("Insert into db error(coin_order_pending). %s, %s, %s, %s, %v", c.GetString("account"), to, coin, amount, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func CancelPendingCoinOrder(c *gin.Context) {
	from := c.GetString("account")
	err := model.MovePendingCoinOrderToCancel(from)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending coin order",
		})
		gl.OutLogger.Error("cancel pending_collect_order error. %s, %v", from, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func GetPendingCoinOrder(c *gin.Context) {
	from := c.GetString("account")
	o, err := model.GetPendingCoinOrder(from)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending coin order",
		})
		gl.OutLogger.Error("get coin_order_pending error. %s, %v", from, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}

func GetCoinOrders(c *gin.Context) {
	from := c.GetString("account")
	count, _ := strconv.Atoi(c.DefaultQuery("count", "5"))
	if count == 0 {
		count = 5
	}
	o, err := model.GetCoinOrders(from, count)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no coin orders",
		})
		gl.OutLogger.Error("get coin_order error. %s, %s, %v", from, c.Query("count"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}

func GetCoinOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "order id error",
		})
		gl.OutLogger.Error("coin order id error. %d, %v", id, err)
		return
	}

	o, err := model.GetCoinOrder(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no order id",
		})
		gl.OutLogger.Error("get coin_order error. %d, %v", id, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}
