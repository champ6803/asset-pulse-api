package handler

import (
	"asset-pulse-api/entities"
	"asset-pulse-api/handler/dto"
	"asset-pulse-api/services"
	"encoding/json"
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

func (h *Handler) GenerateGroupedSoftware(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Query all existing licenses
	licenses, err := h.useCase.GetAllLicenses(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get software licenses: " + err.Error()})
		return
	}

	if len(licenses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no software licenses found"})
		return
	}

	// 2. Convert entities to service input type
	var svcInput []services.SoftwareLicense
	for _, l := range licenses {
		svcInput = append(svcInput, services.SoftwareLicense{
			Name:              l.Name,
			THBPricePerYear:   l.THBPricePerYear,
			UsedByCompanyName: l.UsedByCompanyName,
		})
	}

	// 3. Call Azure OpenAI service
	groupedSvc := h.softwareGroupingService
	grouped, err := groupedSvc.GroupSoftwareLicenses(ctx, svcInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate grouped software: " + err.Error()})
		return
	}

	// 4. Marshal the response back to bytes for upsert
	groupedBytes, err := json.Marshal(grouped)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal grouped software: " + err.Error()})
		return
	}

	// 5. Upsert into DB
	if err := h.useCase.UpsertGroupedSoftware(ctx, groupedBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upsert grouped software: " + err.Error()})
		return
	}

	// 6. Return structured JSON
	c.JSON(http.StatusOK, dto.Response[[]services.GroupedSoftwareResponse]{Data: grouped})
}
