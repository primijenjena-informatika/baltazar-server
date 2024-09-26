package main

import (
	"github.com/primijenjan-informatika/baltazar-server/config"
)

func main() {
	config.InitConfig()
	config.StartServer()
}
