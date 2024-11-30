// prometheusgin/handler.go

package prometheusgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MetricsHandler(reg *MetricRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsData := reg.ExportAll()
		c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte(metricsData))
	}
}
