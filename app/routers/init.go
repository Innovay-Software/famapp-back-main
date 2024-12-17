package routers

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/innovay-software/famapp-main/app/handlers"
// )

// func InitGinEngine(staticMap map[string]string) *gin.Engine {
// 	router := gin.Default()
// 	// gin.Default() already uses gin.Logger
// 	// so calling router.Use(gin.Logger()) is not necessary
// 	// router.Use(gin.Logger())

// 	// Recovery middleware recovers from any panics and writes a custom json response.
// 	router.Use(gin.CustomRecovery(handlers.ApiPanicRecoverHandler))

// 	// Default 404 not found response
// 	router.NoRoute(handlers.Api404Handler)

// 	// Static files
// 	for k, v := range staticMap {
// 		router.Static(k, v)
// 	}

// 	// Web Routes
// 	registerWebV1Routes(router)

// 	// Api Routes
// 	registerApiV2Routes(router)

// 	return router
// }
