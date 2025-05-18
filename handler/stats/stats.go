package stats

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortenLink/services"
	"shortenLink/utils"
	"time"
)

func StatsHandler(c *gin.Context) {
	//获取参数
	code := c.Param("code")
	if !utils.IsValidCode(code) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid short code",
		})
	}

	visits, lastVisit, err := services.StatsService(code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"visits":     visits,
		"last_visit": lastVisit.Format(time.RFC3339),
	})
}
