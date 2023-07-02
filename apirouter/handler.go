package apirouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *APIRouter) indexHandler(c *gin.Context) {
	response := IndexResponse{
		Title:   "_TITLE_",
		Header:  "_HEADER_",
		Block1: "_BLOCK1_",
		Block2: "_BLOCK2_",
	}

	c.HTML(
		http.StatusOK,
		"index.html",
		&response,
	)
}

func (r *APIRouter) pingHandler(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		PingResponse{Message: "pong"},
	)
}
