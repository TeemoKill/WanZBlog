package apirouter

import (
	"database/sql"
	"net/http"

	"github.com/TeemoKill/WanZBlog/config"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type API interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type APIRouter struct {
	Cfg    *config.Config
	router *gin.Engine
	db *sql.DB
}

func New(cfg *config.Config, db *sql.DB) *APIRouter {
	router := &APIRouter{
		Cfg: cfg,
		db: db,
	}
	router.init()

	return router
}

func (r *APIRouter) init() {
	// init gin router
	gin.DisableConsoleColor()
	gin.SetMode(r.Cfg.GinMode)

	ginRouter := gin.New()
	ginRouter.LoadHTMLGlob("template/**/*")
	ginRouter.Static("./static", "static")
	ginRouter.Use(sentrygin.New(sentrygin.Options{}))

	pprof.Register(ginRouter)

	r.router = ginRouter

	// register handlers
	r.router.GET("/", r.indexHandler)
	r.router.GET("/ping", r.pingHandler)

}

func (r *APIRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
