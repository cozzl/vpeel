package main

import (
	"vpeel/api"
	"vpeel/internal/log"
	"vpeel/internal/trans"
	sfu "vpeel/internal/webrtc"
)

func main() {
	log.InitLogger()
	defer log.SyncLogger()

	trans.DefaultManager.Start()
	go trans.DefaultManager.Result()
	go api.Run(":8080")
	go sfu.Run(":8090")
	select {}
}
