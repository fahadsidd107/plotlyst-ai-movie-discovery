package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"plotlyst/server/apiresponse"
	"plotlyst/server/storydna"
	"plotlyst/server/tmdb"
)

type StoryDNAController struct{ service *storydna.Service }

func NewStoryDNAController(service *storydna.Service) *StoryDNAController {
	return &StoryDNAController{service: service}
}

func (h *StoryDNAController) Movie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_movie_id", "A positive numeric movie id is required")
		return
	}
	profile, err := h.service.Profile(c.Request.Context(), id, c.Query("include_spoilers") == "true")
	if err != nil {
		h.failure(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": profile, "source": profile.Provenance})
}

func (h *StoryDNAController) Discover(c *gin.Context) {
	var request storydna.DiscoveryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_story_dna_request", "A valid JSON discovery request is required")
		return
	}
	seed, matches, limitations, err := h.service.Discover(c.Request.Context(), request)
	if err != nil {
		h.failure(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"seed": seed, "matches": matches, "weights": request.Weights}, "provenance": seed.Provenance, "limitations": limitations})
}

func (h *StoryDNAController) failure(c *gin.Context, err error) {
	var upstream *tmdb.APIError
	if errors.As(err, &upstream) {
		status := http.StatusBadGateway
		if upstream.Status == http.StatusNotFound {
			status = http.StatusNotFound
		}
		apiresponse.Error(c, status, "story_dna_source_error", "Unable to assemble Story DNA from the permitted catalogue metadata")
		return
	}
	apiresponse.Error(c, http.StatusBadRequest, "story_dna_unavailable", err.Error())
}
