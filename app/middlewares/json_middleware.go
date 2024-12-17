package middlewares

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/innovay-software/famapp-main/app/services"
// )

// // Processes payload data and convert them to map
// func PayloadValidationMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if err := services.SetPayloadMap(c); err != nil {
// 			panic(err)
// 		}
// 		c.Next()
// 	}
// }

// // Processes url placeholder param data and convert them to map
// func ParamValidationMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if err := services.SetParamMap(c); err != nil {
// 			panic(err)
// 		}
// 		c.Next()
// 	}
// }
