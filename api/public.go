package api

import (
	"iota_dex/gl"
	"iota_dex/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetPairs(c *gin.Context) {
	data, err := model.GetPairs()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.SYSTEM_ERROR,
			"err_msg":  "system error",
		})
		gl.OutLogger.Error("Get pairs from db error. %v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   data,
	})
}

func GetPair(c *gin.Context) {
	coin1 := c.Query("coin1")
	coin2 := c.Query("coin2")
	coin1 = strings.ToUpper(coin1)
	coin2 = strings.ToUpper(coin2)
	if coin1 > coin2 {
		coin1, coin2 = coin2, coin1
	}

	sp, err := model.GetPair(coin1, coin2)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "params error",
		})
		gl.OutLogger.Error("Get price error. %s : %s : %v", coin1, coin2, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   sp,
	})
}

func GetCoin(c *gin.Context) {
	symbol := strings.ToUpper(c.Query("symbol"))

	coin, err := model.GetCoin(symbol)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": gl.PARAMS_ERROR,
			"err_msg":  "params error",
		})
		gl.OutLogger.Error("Get coin info error. %s : %v", symbol, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"data":   coin,
	})
}
