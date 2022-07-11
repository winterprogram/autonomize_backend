package controller

import (
	"io/ioutil"
	"net/http"

	"test/test_app/app/service/correlation"
	"test/test_app/app/service/logger"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func (h HealthController) Status(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)
	log.Info("Inside Health Check Controller.")
	version, err := ioutil.ReadFile("deployedVersion.txt")
	if err != nil {
		log.Info("failed to read file deployedVersion")
	}
	c.JSON(http.StatusOK,
		gin.H{"Version": string(version)})
	return
}
