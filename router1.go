package types

import (
	"os"
	//"techtrainingcamp-AppUpgrade/admin"
	//"techtrainingcamp-AppUpgrade/service"

	"github.com/gin-gonic/gin"
)

func customizeouter(r *gin.Engine) {

	//r.GET("/ping", service.Pong)
	// r.GET("/judge1", service.Hit)
	// r.GET("/judge2", service.HitSQL)
	//r.GET("/judge", service.Judge)
	//r.GET("/count", service.Count)

}

func adminRouter(r *gin.Engine) {
	if os.Getenv("IS_DOCKER") == "1" {
		r.LoadHTMLFiles("/root/public/index.html")
	} else {
		r.LoadHTMLFiles("./public/index.html")
	}
	/*r.GET("/index", admin.GetHTML)
	r.GET("/query_all_rules", admin.QueryAllRules)
	r.GET("/query_rule", admin.QueryRule)
	r.POST("/update_rule", admin.UpdateRule)
	r.POST("/create_rule", admin.CreateRule)
	r.GET("/delete_rule", admin.DeleteRule)
	r.GET("/disable_rule", admin.DisableRule)*/
}
