package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
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

	routerGroup.Handle("GET", "friends", FetchFriendListController)
	routerGroup.Handle("GET", "favorites.json", FetchFavoriteListController)

	engine.Run(":11030")
}

func FetchFriendListController(context *gin.Context) {
	apMsg, errorMessage := FetchRemote(
		"https://raw.githubusercontent.com/liu599/ServiceApi/master/friends.json",
	)
	if errorMessage != "" {
		RespondError(context, http.StatusInternalServerError, errorMessage)
		return
	}
	mk := make(map[string]interface{})
	mk["data"] = apMsg
	Respond(context, http.StatusOK, mk)
}


func FetchFavoriteListController(context *gin.Context) {
	apMsg, errorMessage := FetchRemote(
		"https://raw.githubusercontent.com/liu599/ServiceApi/master/favorites.json",
	)
	if errorMessage != "" {
		RespondError(context, http.StatusInternalServerError, errorMessage)
		return
	}
	mk := make(map[string]interface{})
	mk["data"] = apMsg
	Respond(context, http.StatusOK, mk)
}


func FetchRemote(address string) (appMsg map[string]interface{}, errMsg string) {
	//f, err := os.OpenFile("friends.json", os.O_RDONLY, 0)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", address, nil)
	resp, err := client.Do(req)
	//defer f.Close()
	if err != nil {
		errMsg = fmt.Sprintf("Reading Document : ERROR : %v", err)
		return nil, errMsg
	}
	//data, err := ioutil.ReadAll(f)
	data, err := ioutil.ReadAll(resp.Body)
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