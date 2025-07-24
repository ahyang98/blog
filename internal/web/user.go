package web

import (
	"blog/internal/domain"
	"blog/internal/service"
	ijwt "blog/internal/web/jwt"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	emailPattern    = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	passwordPattern = `^(?=.*\d)(?=.*[a-zA-Z])(?=.*[^\da-zA-Z\s]).{8,20}$`
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	jwtHandler     ijwt.Handler
	svc            service.UserService
}

func NewUserHandler(svc service.UserService, handler ijwt.Handler) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordPattern, regexp.None),
		jwtHandler:     handler,
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/users")
	group.POST("signup", h.SignUp)
	group.POST("login", h.LoginJWT)
	//group.POST("refresh_token", h.RefreshToken)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
		UserName        string `json:"user_name"`
	}
	var req SignUpReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, ErrSystemError)
		return
	}

	if !isEmail {
		ctx.String(http.StatusOK, "email format is not correct")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, ErrSystemError)
		return
	}

	if !isPassword {
		ctx.String(http.StatusOK, "password format is not correct")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "password is not match the confirmed one")
		return
	}

	err = h.svc.SignUp(ctx, &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Username: req.UserName,
	})
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "%s sign up success", req.Email)
	case errors.Is(err, service.ErrDuplicateEmail):
		ctx.String(http.StatusOK, "The email is  already registered")
	default:
		ctx.String(http.StatusOK, ErrSystemError)
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "system error %s", err)
		return
	}
	user, err := h.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	switch {
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "invalid user or password")
	case err == nil:
		err = h.jwtHandler.SetLoginToken(ctx, user.Id)
		if err != nil {
			ctx.String(http.StatusOK, ErrSystemError)
			return
		}
		ctx.JSON(http.StatusOK, "login success")
	default:
		ctx.String(http.StatusOK, ErrSystemError)
	}
}

//func (h *UserHandler) RefreshToken(context *gin.Context) {
//
//}
