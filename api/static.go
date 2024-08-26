package api

import "github.com/gin-gonic/gin"


func staticInit(r *gin.Engine) {
	r.Static("/web", "../web/video_play/")
}

func init() {
	AddRouter(staticInit)
}