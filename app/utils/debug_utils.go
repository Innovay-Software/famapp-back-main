package utils

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func PrintAllRoutes(router *gin.Engine) {
	for _, item := range router.Routes() {
		println("method:", item.Method, "path:", item.Path)
	}
}

func GetCliArgsMap() map[string]string {
	args := os.Args
	argsMap := make(map[string]string)
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			components := strings.Split(arg, "=")
			argsMap[components[0]] = components[1]
		} else {
			argsMap[arg] = "1"
		}
	}
	return argsMap
}
