package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opensourceways/xihe-aicc-finetune/app"
	"github.com/opensourceways/xihe-aicc-finetune/controller"
	"github.com/opensourceways/xihe-aicc-finetune/docs"
)

type Service struct {
	Log      *logrus.Entry
	Port     int
	Timeout  time.Duration
	Finetune app.FinetuneService
}

func StartWebServer(service *Service) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())

	setRouter(r, service)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", service.Port),
		Handler: r,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, service.Timeout)
}

// setRouter init router
func setRouter(engine *gin.Engine, service *Service) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "xihe-aicc-finetune"
	docs.SwaggerInfo.Description = "APIs of xihe finetune"

	v1 := engine.Group(docs.SwaggerInfo.BasePath)
	{
		controller.AddRouterForAICCFinetuneController(
			v1,
			service.Finetune,
		)
	}

	engine.UseRawPath = true
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		logrus.Infof(
			"| %d | %d | %s | %s |",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}
