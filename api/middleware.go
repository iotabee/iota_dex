package api

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"iota_dex/config"
	"iota_dex/gl"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	iotago "github.com/iotaledger/iota.go/v2"
)

func VerifySignature(c *gin.Context) {
	//get user's public key
	sign := c.Query("sign")
	ts := c.Query("ts")
	address := c.Query("addresss")
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
	if timeStamp+config.TokenTime < time.Now().Unix() {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err-code": gl.SIGN_ERROR,
			"err-msg":  "sign expired",
		})
		return
	}

	if len(address) == 42 {
		hash := crypto.Keccak256Hash([]byte(ts))
		err = verifyEthAddress(address, signature, hash.Bytes())
	} else {
		var pk []byte
		pk, err = hex.DecodeString(address)
		if err == nil {
			if err = verifyIotaAddress(pk, signature, []byte(ts)); err == nil {
				ed25519Addr := iotago.AddressFromEd25519PubKey(pk)
				address = ed25519Addr.Bech32(iotago.PrefixMainnet)
			}
		}
	}
	if err != nil {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err-code": gl.SIGN_ERROR,
			"err-msg":  "sign expired",
		})
		gl.OutLogger.Error("User's sign error. %s: %v", address, err)
		return
	}

	c.Set("account", address)
	c.Next()
}

func verifyEthAddress(address string, signature, hashData []byte) error {
	sigPublicKey, err := crypto.SigToPub(hashData, signature)
	if err != nil {
		return errors.New("sign error")
	}
	if address != crypto.PubkeyToAddress(*sigPublicKey).Hex() {
		return errors.New("sign address error")
	}
	return nil
}

func verifyIotaAddress(pubKey, signature, hashData []byte) error {
	if !ed25519.Verify(pubKey, hashData, signature) {
		return errors.New("sign error")
	}
	return nil
}
