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
	address := c.GetString("account")

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

	if err := model.InsertPendingCollectOrder(account, address, coin, amount); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "maybe you have a collect order is pending.",
		})
		gl.OutLogger.Error("Insert into db error(collect_order_pending). %s, %s, %s, %s, %v", account, address, coin, amount, err)
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
	account := c.GetString("account")

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

	if err := model.RetrieveCoin(account, to, coin, amount); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "maybe balance is not enough",
		})
		gl.OutLogger.Error("retriveve coin in db error. %s, %s, %s, %s, %v", account, to, coin, amount, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func CancelPendingCollectOrder(c *gin.Context) {
	address := c.GetString("account")
	err := model.MovePendingCollectOrderToCancel(address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending coin order",
		})
		gl.OutLogger.Error("cancel collect_order_pending error. %s, %v", address, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func GetPendingCollectOrder(c *gin.Context) {
	address := c.GetString("account")
	o, err := model.GetPendingCollectOrder(address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending collect order",
		})
		gl.OutLogger.Error("get collect_order_pending error. %s, %v", address, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}

func GetCoinOrders(c *gin.Context) {
	address := c.GetString("account")
	count, _ := strconv.Atoi(c.DefaultQuery("count", "5"))
	if count == 0 {
		count = 5
	}
	if count > 100 {
		count = 100
	}
	o, err := model.GetCoinOrders(address, count)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no coin orders",
		})
		gl.OutLogger.Error("get coin_order error. %s, %s, %v", address, c.Query("count"), err)
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
