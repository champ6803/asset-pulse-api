package handler

import (
	"asset-pulse-api/entities"
	"asset-pulse-api/handler/dto"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllLicenses(c *gin.Context) {
	licenses, err := h.useCase.GetAllLicenses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []dto.InsertLicenseRequest
	for _, l := range licenses {
		resp = append(resp, dto.InsertLicenseRequest{
			ID:                int(l.ID),
			Name:              l.Name,
			ThbPricePerYear:   l.THBPricePerYear,
			UsedByCompanyName: l.UsedByCompanyName,
		})
	}

	c.JSON(http.StatusOK, dto.Response[[]dto.InsertLicenseRequest]{Data: resp})
}

func (h *Handler) InsertLicense(c *gin.Context) {
	var req dto.InsertLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	license := entities.SoftwareLicense{
		Name:              req.Name,
		THBPricePerYear:   req.ThbPricePerYear,
		UsedByCompanyName: req.UsedByCompanyName,
	}

	if err := h.useCase.InsertLicense(c.Request.Context(), &license); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.Response[dto.InsertLicenseRequest]{Data: req})
}

func (h *Handler) DeleteLicense(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.useCase.DeleteLicense(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "deleted"})
}

func (h *Handler) GetGroupedSoftware(c *gin.Context) {
	groups, err := h.useCase.GetGroupedSoftware(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response[[]entities.CurrentGroupedSoftware]{Data: groups})
}

func (h *Handler) UpsertGroupedSoftware(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	if err := h.useCase.UpsertGroupedSoftware(c.Request.Context(), data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

type GroupGenerationRequest struct {
	// Optionally you can add flags here in the future
}

// GenerateGroupedSoftware is deprecated - grouping is now done via category-based approach in GetSimilarSoftwareClusters
// This endpoint is kept for backward compatibility but now returns a deprecation notice
func (h *Handler) GenerateGroupedSoftware(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "This endpoint is deprecated. Please use /api/v1/similar-software/clusters instead which uses category-based grouping via CatalogSearchService.",
		"message": "Grouping is now handled automatically by category from catalog search results",
	})
}
