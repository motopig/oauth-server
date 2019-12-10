package main

import (
	"errors"
	"github.com/go-session/session"

	//"fmt"
	"github.com/motopig/oauth-server/app/model"
	"github.com/motopig/oauth-server/app/router"
	"github.com/motopig/oauth-server/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/gin-server"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	oredis "gopkg.in/go-oauth2/redis.v3"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"log"
	"net/http"
)

func main() {
	config.InitConfig()
	model.InitMysql()
	manager := manage.NewDefaultManager()
	// use redis token store
	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr: viper.GetString("redis.host"),
		DB:   viper.GetInt("redis.database"),
	}))

	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(viper.GetString("jwt.key")), jwt.SigningMethodHS512))

	// client store
	clientStore := store.NewClientStore()
	clients := model.GetAllClients()
	if len(clients) > 0 {
		for _, v := range clients {
			_ = clientStore.Set(v.ClientId, &models.Client{
				ID:     v.ClientId,
				Secret: v.ClientSecret,
				Domain: v.ClientDoamin,
			})
		}
	} else {
		log.Fatal("NO CLIENTS FOUND!!")
	}

	manager.MapClientStorage(clientStore)

	// Initialize the oauth2 service
	ginserver.InitServer(manager)
	ginserver.SetAllowGetAccessRequest(true)
	ginserver.SetClientInfoHandler(server.ClientFormHandler)
	ginserver.SetUserAuthorizationHandler(userAuthorizeHandler) // http://localhost:8080/api/oauth2/authorize?client_id=007&response_type=code&scope=read&state=STATE&redirect_uri=http%3A%2F%2Flocalhost/api/auth
	// http://localhost:8080/api/oauth2/token?client_id=007&client_secret=Fi9GXy6X77dEeZ8t&grant_type=authorization_code&code=JFC0ZWISNYIWVPPVV7FAWA&response_type=token&redirect_uri=http%3A%2F%2Flocalhost/api/auth
	// http://localhost:8080/api/oauth2/token?client_id=007&client_secret=Fi9GXy6X77dEeZ8t&grant_type=refresh_token&refresh_token=JFC0ZWISNYIWVPPVV7FAWA&response_type=token&redirect_uri=http%3A%2F%2Flocalhost/api/auth
	// redirect_uri 要和 clientStore 中保存的domain一致
	port := viper.GetString("system.port")
	host := viper.GetString("system.host")
	r := gin.Default()
	router.Load(r)

	if err := r.Run(host + ":" + port); err != nil {
		log.Fatal(err)
	}
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userId string, err error) {
	sessiond, err := session.Start(nil, w, r)
	if err != nil {
		return "", errors.New("User Not Login!")
	}

	var ok bool
	var user_id interface{}
	if user_id, ok = sessiond.Get("LoggedInUserID"); !ok {
		return "", errors.New("User Not Login!")
	}
	userId = user_id.(string)
	return userId, nil
}
