package main

import (
	"vpeel/api"
	"vpeel/internal/log"
	"vpeel/internal/trans"
)

func main() {
	log.InitLogger()
	defer log.SyncLogger()

	trans.DefaultManager.Start()
	go trans.DefaultManager.Result()

	api.Run()
}
