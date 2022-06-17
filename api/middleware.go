package api

import (
	"iota_dex/config"
	"iota_dex/gl"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

func VerifySignature(c *gin.Context) {
	//get user's public key
	sign := c.Query("sign")
	ts := c.Query("ts")
	hash := crypto.Keccak256Hash([]byte(ts))
	signature, err := hexutil.Decode(sign)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err-code": gl.PARAMS_ERROR,
			"err-msg":  "invalid sign",
		})
		return
	}
	timeStamp, _ := strconv.ParseInt(ts, 10, 64)
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if (err != nil) || (timeStamp+config.TokenTime < time.Now().Unix()) {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err-code": gl.SIGN_ERROR,
			"err-msg":  "sign error or expired",
		})
		return
	}
	c.Set("account", hexutil.Encode(sigPublicKey))
	c.Next()
}
