package main

import (
	"iota_dex/api"
	"iota_dex/gl"
	"iota_dex/model"
)

func main() {
	gl.CreateLogFiles()

	model.ConnectToMysql()

	api.StartHttpServer()

}
