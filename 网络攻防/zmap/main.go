package main

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"net/http"
	"zmap-frontend/utili"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	go func() {
		r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * 3600,
		}))
		r.POST("/", func(c *gin.Context) {
			var requestData map[string]interface{}

			if err := c.BindJSON(&requestData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}
			fmt.Println(requestData) // 打印接收到的数据
			res := make(map[string]string)
			res = utili.Handle(requestData)
			//var res = make(map[string]interface{})
			//if num, ok := requestData["scannum"].(int); ok {
			//	if num > 1 {
			//		for i := 0; i < num; i++ {
			//			manyres := utili.Handle(requestData)
			//			res[strconv.Itoa(i)] = manyres
			//		}
			//	} else {
			//		oneres := utili.Handle(requestData)
			//		res["1"] = oneres
			//	}
			//}

			fmt.Println(res)
			c.JSON(http.StatusOK, res)
		})

		// 启动 Gin 服务器
		if err := r.Run(":7899"); err != nil {
			fmt.Printf("Failed to run server: %v\n", err)
		}
	}()
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "zmap-frontend",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
