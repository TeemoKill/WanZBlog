package apirouter

import (
	"github.com/TeemoKill/WanZBlog/apirouter/requests"
	"github.com/TeemoKill/WanZBlog/log"
	"net/http"

	"github.com/TeemoKill/WanZBlog/datamodel"
	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

type RegisterResponse struct {
}

func (r *APIRouter) registerPageHandler(c *gin.Context) {
	logger := log.CurrentModuleLogger()
	logger.Info("registerPageHandler")

	response := RegisterResponse{}

	c.HTML(
		http.StatusOK,
		"register.html",
		&response,
	)
}

type RegisterResultResponse struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`

	UserEmail string `json:"user_email"`
	Token     string `json:"token"`
}

func (r *APIRouter) registerHandler(c *gin.Context) {
	logger := log.CurrentModuleLogger()
	logger.Info("registerHandler")

	request := &requests.RegisterRequest{}
	if err := c.ShouldBind(request); err != nil {
		logger.WithError(err).
			Errorf("gin context bind parameter error")
		c.JSON(http.StatusOK,
			RegisterResultResponse{
				Code:    1,
				Message: "gin context bind parameter error",
			},
		)
		return
	}
	logger.WithField("email", request.Email).
		WithField("username", request.Username).
		Infof("request parameters")

	// check if the email is registered
	user := &datamodel.User{
		Email: request.Email,
	}
	var userCount int64
	if err := r.db.Model(user).Where("email=?", user.Email).Count(&userCount).Error; err != nil {
		logger.WithError(err).Warn("fetch user info error")
		c.JSON(http.StatusOK,
			RegisterResultResponse{
				Code:    1,
				Message: "fetch user info error",
			},
		)
		return
	}
	if userCount > 0 {
		logger.WithField("email", request.Email).
			Info("email has already registered")
		c.JSON(http.StatusOK,
			RegisterResultResponse{
				Code:    1,
				Message: "email has already registered",
			},
		)
	}

	// TODO: verify email address

	encryptedPw, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithError(err).Warn("encrypt password error")
		c.JSON(http.StatusOK,
			RegisterResultResponse{
				Code:    1,
				Message: "encrypt password error",
			},
		)
	}

	newUser := &datamodel.User{
		Email:    request.Email,
		Password: string(encryptedPw),
		Username: request.Username,
	}
	if err := newUser.Create(r.db); err != nil {
		logger.WithError(err).Warn("create user error")
		c.JSON(http.StatusOK,
			RegisterResultResponse{
				Code:    1,
				Message: "create user error",
			},
		)
	}

	token, err := GenerateLoginToken(newUser.UUID)
	if err != nil {
		logger.WithError(err).Error("GenerateLoginToken error")
		c.JSON(http.StatusOK,
			RegisterResultResponse{
				Code:    1,
				Message: err.Error(),
			},
		)
	}

	logger.WithField("token", token).
		WithField("email", request.Email).
		WithField("username", request.Username).
		Info("register success")

	c.JSON(http.StatusOK,
		RegisterResultResponse{
			Code:      0,
			UserEmail: newUser.Email,
			Token:     token,
			Message:   "register success",
		},
	)

}
