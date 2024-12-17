package routers

// import (
// 	"net/http"
// 	"os"
// 	"text/template"

// 	"github.com/innovay-software/famapp-main/app/handlers"
// 	"github.com/innovay-software/famapp-main/app/utils"

// 	"github.com/gin-contrib/multitemplate"
// 	"github.com/gin-gonic/gin"
// )

// // Register functions to expose to html templates
// func createFuncMap() template.FuncMap {
// 	return template.FuncMap{
// 		"cdnRoot": func() string {
// 			return os.Getenv("CDN_ROOT")
// 		},
// 	}
// }

// // Register WEB routes
// func registerWebV1Routes(router *gin.Engine) {

// 	viewDir := utils.GetRootAbsPath("views")
// 	funcMap := createFuncMap()

// 	r := multitemplate.NewRenderer()
// 	r.AddFromFilesFuncs("login", funcMap, viewDir+"/layouts/layout.html", viewDir+"/pages/login.html")
// 	router.HTMLRender = r

// 	router.GET("/", func(c *gin.Context) {
// 		c.Redirect(http.StatusMovedPermanently, os.Getenv("DEV_HOME"))
// 	})
// 	router.GET("/login", handlers.LoginPage)
// }
