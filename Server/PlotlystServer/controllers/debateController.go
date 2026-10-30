package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"plotlyst/server/apiresponse"
	"plotlyst/server/debate"
)

type DebateController struct{ service *debate.Service }

func NewDebateController(service *debate.Service) *DebateController {
	return &DebateController{service: service}
}
func (h *DebateController) Movie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_movie_id", "A positive movie id is required")
		return
	}
	bucket := c.DefaultQuery("rating_bucket", "all")
	if bucket != "all" && bucket != "high" && bucket != "middle" && bucket != "low" {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_rating_bucket", "Use all, high, middle, or low")
		return
	}
	summary, err := h.service.Get(c.Request.Context(), id, bucket)
	if err != nil {
		apiresponse.Error(c, http.StatusInternalServerError, "debate_unavailable", "The audience debate could not be assembled")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": summary, "source": gin.H{"provider": "Plotlyst community reviews", "generated": summary.AnalysisMethod != "none", "analysis_method": summary.AnalysisMethod}})
}
