package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	engine.Use(gin.Logger())

	engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "User", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Real-Ip"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	}))

	routerGroup := engine.Group("v1/apps")

	routerGroup.Handle("GET", "friends", ReadFileController)

	engine.Run(":11030")
}

func ReadFileController(context *gin.Context) {
	apMsg, errorMessage := ReadFile()
	if errorMessage != "" {
		RespondError(context, http.StatusInternalServerError, errorMessage)
		return
	}
	mk := make(map[string]interface{})
	mk["data"] = apMsg
	Respond(context, http.StatusOK, mk)
}

func ReadFile() (appMsg map[string]interface{}, errMsg string) {
	f, err := os.OpenFile("friends.json", os.O_RDONLY, 0)
	defer f.Close()
	if err != nil {
		errMsg = fmt.Sprintf("Reading Document : ERROR : %v", err)
		return nil, errMsg
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		errMsg = fmt.Sprintf("Reading Document : ERROR : %v", err)
		return nil, errMsg
	}
	var friendList map[string]interface{}
	if err = json.Unmarshal([]byte(data), &friendList); err != nil {
		errMsg = fmt.Sprintf("JSON Unmarshal : ERROR : %v", err)
		return nil, errMsg
	}
	return friendList, ""
}