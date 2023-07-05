package apirouter

import (
	"github.com/TeemoKill/WanZBlog/datamodel"
	"github.com/TeemoKill/WanZBlog/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserProfileResponse struct {
	Username        string `json:"username"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_image_url"`
}

func (r *APIRouter) userProfileHandler(c *gin.Context) {
	logger := log.CurrentModuleLogger()
	logger.Info("userProfileHandler")

	/*
		request := &requests.UserProfileRequest{}
		if err := c.ShouldBind(request); err != nil {
			logger.WithError(err).
				Errorf("gin context bind parameter error")
			c.HTML(http.StatusNotFound, "404.html", NotFoundResponse{})
			return
		}
	*/
	logger.WithField("user_uuid", c.Param("user_uuid")).
		Infof("request parameters")

	user := &datamodel.User{}
	user.UUID = c.Param("user_uuid")
	if err := user.LoadByUUID(r.db); err != nil {
		logger.WithError(err).Warn("fetch user info error")
		c.HTML(http.StatusNotFound, "404.html", NotFoundResponse{})
		return
	}
	logger.WithField("user_id", user.ID).
		WithField("user_uuid", user.UUID).
		WithField("username", user.Username).
		WithField("email", user.Email).
		Info("user info")

	response := UserProfileResponse{
		Username: user.Username,
		Bio:      user.Bio,
	}

	c.HTML(
		http.StatusOK,
		"profile.html",
		&response,
	)
}
