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
	postDaoGorm := dao.NewPostDaoGorm(db)
	postRepository := repository.NewPostRepository(postDaoGorm, userRepository)
	postService := service.NewPostService(postRepository, loggerV1)
	postHandler := web.NewPostHandler(postService, loggerV1)
	commentGormDao := dao.NewCommentGormDao(db)
	commentRepo := repository.NewCommentRepo(commentGormDao, userRepository)
	commentService := service.NewCommentService(commentRepo, postService, loggerV1)
	commentHandler := web.NewCommentHandler(commentService, loggerV1)
	webServer := ioc.InitWebServer(ginMiddlewares, userHandler, postHandler, commentHandler)
	loggerV1.Info("webserver initialed.")
	err := webServer.Run(":8081")
	if err != nil {
		panic(err)
	}
}

func initViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "config file path")
	pflag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config fail %s", err))
	}
	fmt.Println(viper.Get("test.key"))
}
