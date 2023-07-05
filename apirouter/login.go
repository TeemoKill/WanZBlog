package apirouter

import (
	"github.com/TeemoKill/WanZBlog/apirouter/requests"
	"github.com/TeemoKill/WanZBlog/datamodel"
	"github.com/TeemoKill/WanZBlog/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginResponse struct {
}

func (r *APIRouter) loginPageHandler(c *gin.Context) {
	logger := log.CurrentModuleLogger()
	logger.Info("loginPageHandler")

	response := LoginResponse{}

	c.HTML(
		http.StatusOK,
		"login.html",
		&response,
	)
}

type LoginResultResponse struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`

	UserUUID  string `json:"user_uuid"`
	UserEmail string `json:"user_email"`
	Token     string `json:"token"`
}

func (r *APIRouter) loginHandler(c *gin.Context) {
	logger := log.CurrentModuleLogger()
	logger.Info("loginHandler")

	request := &requests.LoginRequest{}
	if err := c.ShouldBind(request); err != nil {
		logger.WithError(err).
			Errorf("gin context bind parameter error")
		c.JSON(http.StatusOK,
			LoginResultResponse{
				Code:    1,
				Message: "gin context bind parameter error",
			},
		)
		return
	}
	logger.WithField("email", request.Email).
		Infof("request parameters")

	// Check if the email is registered
	user := &datamodel.User{
		Email: request.Email,
	}
	if err := user.LoadByEmail(r.db); err != nil {
		logger.WithError(err).Warn("fetch user info error")
		c.JSON(http.StatusOK,
			LoginResultResponse{
				Code:    1,
				Message: "fetch user info error",
			},
		)
		return
	}

	// Check if the password matches
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		logger.WithError(err).Warn("incorrect password")
		c.JSON(http.StatusOK,
			LoginResultResponse{
				Code:    1,
				Message: "incorrect password",
			},
		)
		return
	}

	token, err := GenerateLoginToken(user.UUID)
	if err != nil {
		logger.WithError(err).Error("GenerateLoginToken error")
		c.JSON(http.StatusOK,
			LoginResultResponse{
				Code:    1,
				Message: err.Error(),
			},
		)
	}

	logger.WithField("token", token).
		WithField("email", request.Email).
		Info("login success")

	c.JSON(http.StatusOK,
		LoginResultResponse{
			Code:      0,
			UserUUID:  user.UUID,
			UserEmail: user.Email,
			Token:     token,
			Message:   "login success",
		},
	)
}
