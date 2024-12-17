package main

import (
	// "encoding/json"
	"path/filepath"
	"runtime"

	"github.com/innovay-software/famapp-main/app"

	// "github.com/innovay-software/famapp-main/app/repositories"
	// "github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	projDir := filepath.Dir(b)
	argsMap := utils.GetCliArgsMap()

	switch argsMap["-env"] {
	case "prod":
		utils.LogSuccess("Starting Prod Server:", projDir)
		ginInstance, port := app.InitApiProdServer(projDir)
		ginInstance.Run(":" + port)
	case "local":
		utils.LogSuccess("Start Local Server:", projDir)
		ginInstance, port := app.InitApiLocalServer(projDir)
		utils.PrintAllRoutes(ginInstance)
		ginInstance.Run(":" + port)
	default:
		utils.LogError("Missing -env cli argument:", argsMap["-env"])
	}

	utils.LogError("Gin server terminated...")
}

// func test() {
// 	jsonUser := `
// 	{
// 		"userData": {
// 			"id": 1,
// 			"name": "test",
// 			"createdAt": "2021-04-04T11:33:33Z"
// 		}
// 	}
// 	`

// 	mapUser := map[string]any{
// 		"userData": map[string]any{
// 			"id": 1,
// 			"name": "test",
// 			"createdAt": "2021-04-04T11:33:33Z",
// 		},
// 	}
// 	// var buf bytes.Buffer
// 	// enc := gob.NewEncoder(&buf)
// 	// err := enc.Encode(jsonUser)
// 	// if err != nil {
// 	// 	utils.LogError("Encode err:", err.Error())
// 	// }

// 	// type embededUser struct {
// 	// 	User       *models.User         `json:"userData"`
// 	// }
// 	// var eUser embededUser
// 	var res dto.LoginResponse
// 	err2 := json.Unmarshal([]byte(jsonUser), &res)
// 	if err2 != nil {
// 		utils.LogError("Encode err:", err2.Error())
// 	}
// 	utils.LogError("User = ", res.User)
// 	// Output:
// 	// `
// 	// User =  {{1 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC} 00000000-0000-0000-0000-000000000000  test  <nil>     <nil> <nil> map[] 0}
// 	// `

// 	bytes, err := json.Marshal(mapUser)
// 	if err != nil {
// 		utils.LogError("marshal map error:", err)
// 	}
// 	var res2 dto.LoginResponse
// 	err3 := json.Unmarshal(bytes, &res2)
// 	if err3 != nil {
// 		utils.LogError("Encode err:", err3.Error())
// 	}
// 	utils.LogError("User2 = ", res2.User)

// }
