package api

import (
	"iota_dex/config"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/triplefi/go-logger/logger"
)

func StartHttpServer() {
	router := InitRouter()
	router.Run(":" + strconv.Itoa(config.HttpPort))
}

// InitRouter init the router
func InitRouter() *gin.Engine {
	if err := os.MkdirAll("./logs/http", os.ModePerm); err != nil {
		log.Panic("Create dir './logs/http' error. " + err.Error())
	}
	GinLogger, err := logger.New("logs/http/gin.log", 2, 100*1024*1024, 10)
	if err != nil {
		log.Panic("Create GinLogger file error. " + err.Error())
	}

	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: GinLogger}), gin.Recovery())
	router.SetTrustedProxies(nil)

	api := router.Group("/api")
	{
		api.GET("/pairs", GetPairs)
		api.GET("/price", GetPrice)
		api.GET("/balance", GetBalance)
	}

	order := api.Group("/order")
	order.Use(VerifySignature)
	{
		order.GET("/swap", SwapOrder)
		order.GET("/swap/pending", GetPendingSwapOrder)
		order.GET("/swap/cancel", CancelPendingSwapOrder)
		order.GET("/swap/list", GetSwapOrders)
		order.GET("/swap/info", GetSwapOrder)

		order.GET("/coin/collect", CollectCoinOrder)
		order.GET("/coin/retrieve", RetrieveCoinOrder)
		order.GET("/coin/pending", GetPendingCoinOrder)
		order.GET("/coin/cancel", CancelPendingCoinOrder)
		order.GET("/coin/list", GetCoinOrders)
		order.GET("/coin/info", GetCoinOrder)
	}

	return router
}
