package apirouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexResponse struct {
	Title  string `json:"title"`
	Header string `json:"header"`
	Block1 string `json:"block1"`
	Block2 string `json:"block2"`
}

type PingResponse struct {
	Message string `json:"message"`
}

func (r *APIRouter) indexHandler(c *gin.Context) {
	response := IndexResponse{
		Title:  "_TITLE_",
		Header: "_HEADER_",
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
