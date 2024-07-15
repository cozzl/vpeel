package api

import "github.com/gin-gonic/gin"

func checkRouterInit(r *gin.Engine) {

	r.GET("/check", check)

}
func check(c *gin.Context) {
	resp := basicResponse{Code: 200, Meeeage: "success"}

	c.JSON(200,resp)
	// c.JSON(200, gin.H{
	// 	"message": "success",
	// })
}

func init() {
	AddRouter(checkRouterInit)
}
