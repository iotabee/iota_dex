package main

import (
	"iota_dex/api"
	"iota_dex/config"
	"iota_dex/daemon"
	"iota_dex/gl"
	"iota_dex/model"
	"iota_dex/test"
)

func main() {
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
