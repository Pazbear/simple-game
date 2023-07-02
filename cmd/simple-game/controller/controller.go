package controller

import (
	"fmt"
	"net/http"
	"whatisthissong/cmd/whatisthissong/config"
	"whatisthissong/cmd/whatisthissong/docs"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controller struct {
	AppConfig config.Config
}

func NewController() (*Controller, error) {
	return &Controller{}, nil
}

// @Summary     상태 체크
// @Description	현재 서버 상태 체크
// @Tags        common
// @Router      /api/v1/healthcheck [get]
// @Success     200
func (c *Controller) HealthCheck(ginctx *gin.Context) {
	ginctx.JSON(http.StatusOK, nil)
}

func (c *Controller) NewRouter() *gin.Engine {

	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/healthcheck", c.HealthCheck)

		}
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", c.AppConfig.Address, c.AppConfig.Port)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
