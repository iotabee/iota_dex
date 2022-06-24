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
	go router.Run(":" + strconv.Itoa(config.HttpPort))
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

	api := gin.New()
	api.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: GinLogger}), gin.Recovery())
	api.SetTrustedProxies(nil)

	public := api.Group("/public")
	{
		public.GET("/pairs", GetPairs)
		public.GET("/price", GetPrice)
		public.GET("/balance", GetBalance)
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
		order.GET("/coin/pending", GetPendingCollectOrder)
		order.GET("/coin/list", GetCoinOrders)
		order.GET("/coin/info", GetCoinOrder)

		order.GET("/liquidity/add", LiquidityAddOrder)
		order.GET("/liquidity/remove", LiquidityRemoveOrder)
		order.GET("/liquidity/pending", GetPendingLiquidityAddOrder)
		order.GET("/liquidity/cancel", CancelPendingLiquidityAddOrder)
		order.GET("/liquidity/list", GetLiquidityOrders)
		order.GET("/liquidity/info", GetLiquidityOrder)
	}

	return api
}
