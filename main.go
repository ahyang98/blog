package main

import (
	"blog/internal/repository"
	"blog/internal/repository/dao"
	"blog/internal/service"
	"blog/internal/web"
	"blog/internal/web/jwt"
	"blog/ioc"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperV1()
	cache := ioc.InitCache()
	jwtHandler := jwt.NewJWTHandler(cache)
	loggerV1 := ioc.InitialLogger()
	ginMiddlewares := ioc.InitGinMiddlewares(cache, jwtHandler)
	db := ioc.InitDB(loggerV1)
	userDao := dao.NewGormUserDao(db)
	userRepository := repository.NewUserRepository(userDao)
	userService := service.NewUserService(userRepository, loggerV1)
	userHandler := web.NewUserHandler(userService, jwtHandler)
	webServer := ioc.InitWebServer(ginMiddlewares, userHandler)
	loggerV1.Info("webserver initialed.")
	err := webServer.Run(":8081")
	if err != nil {
		panic(err)
	}
}

func initViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "config file path")
	pflag.Parse()
	//viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	//viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config fail %s", err))
	}
	fmt.Println(viper.Get("test.key"))
}
