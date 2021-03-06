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

func SwapOrder(c *gin.Context) {
	source := c.Query("source")
	target := c.Query("target")
	to := c.Query("to")
	amount := c.Query("amount")
	min_amount := c.Query("min_amount")
	_, b1 := new(big.Int).SetString(amount, 10)
	_, b2 := new(big.Int).SetString(min_amount, 10)

	coin1 := strings.ToUpper(source)
	coin2 := strings.ToUpper(target)
	if coin1 > coin2 {
		coin1, coin2 = coin2, coin1
	}
	_, err := model.GetPair(coin1, coin2)
	if err != nil || !b1 || !b2 || len(to) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "params error",
		})
		gl.OutLogger.Error("Get price(%s:%s) error when swap_order, or params error(%s,%s,%s). %v", coin1, coin2, amount, min_amount, to, err)
		return
	}

	if err := model.InsertPendingSwapOrder(c.GetString("account"), source, amount, to, target, min_amount); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "maybe you have a swap order is pending",
		})
		gl.OutLogger.Error("Insert into db error(swap_order_pending). %s, %v", c.GetString("account"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func GetPendingSwapOrder(c *gin.Context) {
	account := c.GetString("account")
	o, err := model.GetPendingSwapOrder(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending swap order",
		})
		gl.OutLogger.Error("get pending_swap_order error. %s, %v", account, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}

func CancelPendingSwapOrder(c *gin.Context) {
	account := c.GetString("account")
	err := model.MovePendingSwapOrderToCancel(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no pending swap order",
		})
		gl.OutLogger.Error("cancel swap_order_pending error. %s, %v", account, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func GetSwapOrders(c *gin.Context) {
	account := c.GetString("account")
	count, err := strconv.Atoi(c.DefaultQuery("count", "5"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "params error",
		})
		gl.OutLogger.Error("param error when get swap orders. %s : %v", c.Query("count"), err)
		return
	}
	if count > config.MaxQueryCount {
		count = config.MaxQueryCount
	}

	o, err := model.GetSwapOrders(account, count)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no swap orders",
		})
		gl.OutLogger.Error("get swap_order error. %s, %s, %v", account, c.Query("count"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}

func GetSwapOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "order id error",
		})
		gl.OutLogger.Error("swap order id error. %d, %v", id, err)
		return
	}

	o, err := model.GetSwapOrder(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "have no this order",
		})
		gl.OutLogger.Error("get swap_order error. %d, %v", id, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   o,
	})
}
