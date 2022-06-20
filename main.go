package main

import (
	"fmt"
	"iota_dex/api"
	"iota_dex/config"
	"iota_dex/daemon"
	"iota_dex/gl"
	"iota_dex/model"
	"iota_dex/test"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	privateKey, _ := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	data := []byte("1655714635")       //
	hash := crypto.Keccak256Hash(data) //0xb2572721d4d9b79d1f7678b68da3265dfc42dbae0ce0603cc17e3af26f5deb51
	signature, _ := crypto.Sign(hash.Bytes(), privateKey)
	sign := hexutil.Encode(signature)
	fmt.Println(sign)
	return
	if config.Env == "product" {
		daemon.Background("./out.log", true)
	}

	gl.CreateLogFiles()

	model.ConnectToMysql()

	api.StartHttpServer()

	if config.Env == "dev" {
		test.RunTest()
	}

	daemon.WaitForKill()
}
