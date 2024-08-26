package api

import "github.com/gin-gonic/gin"


func check(c *gin.Context) {
	resp := basicResponseS{Code: 200, Meeeage: "success"}
	c.JSON(200,resp)
}

func checkRouterInit(r *gin.Engine) {
	r.GET("/check", check)
}


func init() {
	AddRouter(checkRouterInit)
}
