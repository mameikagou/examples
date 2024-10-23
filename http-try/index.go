package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/echo", func(ctx *gin.Context) {
		var jsonData map[string]interface{}

		if err := ctx.BindJSON(&jsonData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		method := ctx.Request.Method
		url := ctx.Request.URL.String()

		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.JSON(http.StatusOK, gin.H{
			"method": method,
			"url":    url,
			"body":   jsonData,
		})
	})

	r.Run(":8080")
}

// curl -X POST http://localhost:8080/echo -H "Content-Type: application/json" -d '{"name": "GitHub Copilot"}'

// 返回值:
// {
// 	"method": "POST",
// 	"url": "/echo",
// 	"data": {
// 	  "name": "GitHub Copilot"
// 	}
// }
