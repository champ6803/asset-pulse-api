package handler

import (
	"asset-pulse-api/entities"
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSimilarSoftwareClusters returns consolidation-like snapshot data to power the Similar Software UI
func (h *Handler) GetSimilarSoftwareClusters(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	companyCode := c.Query("company_code")
	selectedSubsidiaries := c.QueryArray("subsidiaries") // Support multiple subsidiaries filter
	_ = c.Query("app_name")                              // Support app name search (reserved for future use)

	// Get consolidation opportunities
	data, err := h.dbRepo.GetConsolidationOpportunities(ctx, companyCode)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Get similar software clusters error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Enrich with subsidiaries data from license_assignments
	if data != nil && len(*data) > 0 {
		appIDs := make([]int64, 0)
		for _, opp := range *data {
			if opp.AppID != nil {
				appIDs = append(appIDs, *opp.AppID)
			}
		}

		subsidiariesMap, err := h.dbRepo.GetAppsWithSubsidiaries(ctx, appIDs)
		if err == nil && len(subsidiariesMap) > 0 {
			// Enrich each opportunity with subsidiaries from license_assignments
			// Store in a custom response structure
			type EnrichedOpportunity struct {
				entities.GroupConsolidationOpp
				Subsidiaries []string `json:"subsidiaries"`
			}

			enriched := make([]EnrichedOpportunity, 0, len(*data))
			for _, opp := range *data {
				subsidiaries := []string{}
				if opp.AppID != nil {
					if subs, ok := subsidiariesMap[*opp.AppID]; ok {
						subsidiaries = subs
					}
				}
				enriched = append(enriched, EnrichedOpportunity{
					GroupConsolidationOpp: opp,
					Subsidiaries:          subsidiaries,
				})
			}

			output := transformer.SuccessResponse(http.StatusOK, enriched)
			c.JSON(http.StatusOK, output)
			return
		}
	}

	_ = selectedSubsidiaries // Reserved for future DB-level filtering

	output := transformer.SuccessResponse(http.StatusOK, data)
	c.JSON(http.StatusOK, output)
}

type pricingTier struct {
	ThresholdQty int     `json:"threshold_qty"`
	UnitPrice    float64 `json:"unit_price"`
}

type vendorPricing struct {
	VendorID          string        `json:"vendor_id"`
	VendorName        string        `json:"vendor_name"`
	FeatureClusterKey string        `json:"feature_cluster_key"`
	PricingMode       string        `json:"pricing_mode"`   // piecewise|progressive
	BillingPeriod     string        `json:"billing_period"` // monthly|yearly
	Currency          string        `json:"currency"`
	Tiers             []pricingTier `json:"tiers"`
}

// GetVendorPricingTiers returns mock vendor pricing tiers per cluster/vendor
func (h *Handler) GetVendorPricingTiers(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	clusterKey := c.Query("cluster_key")
	vendorID := c.Query("vendor_id")

	// Simple mocked tiers; replace with DB pricing tables when available
	tiers := []vendorPricing{}
	if clusterKey == "design_tools" || clusterKey == "" {
		tiers = append(tiers, vendorPricing{
			VendorID:          "figma",
			VendorName:        "Figma Inc.",
			FeatureClusterKey: "design_tools",
			PricingMode:       "piecewise",
			BillingPeriod:     "monthly",
			Currency:          "THB",
			Tiers:             []pricingTier{{1, 600}, {50, 500}, {200, 400}},
		})
		tiers = append(tiers, vendorPricing{
			VendorID:          "adobe",
			VendorName:        "Adobe",
			FeatureClusterKey: "design_tools",
			PricingMode:       "progressive",
			BillingPeriod:     "yearly",
			Currency:          "THB",
			Tiers:             []pricingTier{{1, 12000}, {100, 10000}},
		})
	}
	if clusterKey == "project_mgmt" || clusterKey == "" {
		tiers = append(tiers, vendorPricing{
			VendorID:          "atlassian",
			VendorName:        "Atlassian",
			FeatureClusterKey: "project_mgmt",
			PricingMode:       "piecewise",
			BillingPeriod:     "monthly",
			Currency:          "THB",
			Tiers:             []pricingTier{{1, 450}, {50, 350}, {200, 300}},
		})
	}

	if vendorID != "" {
		filtered := make([]vendorPricing, 0, 1)
		for _, v := range tiers {
			if v.VendorID == vendorID {
				filtered = append(filtered, v)
			}
		}
		tiers = filtered
	}

	output := transformer.SuccessResponse(http.StatusOK, tiers)
	c.JSON(http.StatusOK, output)
}

type savingSimulationRequest struct {
	ClusterKey                  string   `json:"cluster_key" binding:"required"`
	TargetVendorID              string   `json:"target_vendor_id" binding:"required"`
	TrainingCostPerUser         *float64 `json:"training_cost_per_user"`
	MigrationFlatCost           *float64 `json:"migration_flat_cost"`
	EarlyTerminationPenaltyRate *float64 `json:"early_termination_penalty_rate"`
}

type savingSimulationResponse struct {
	CurrentCostYear    float64 `json:"current_cost_year"`
	ProposedLicenses   float64 `json:"proposed_licenses"`
	SwitchingCost      float64 `json:"switching_cost"`
	ProposedTotal      float64 `json:"proposed_total"`
	PotentialSaving    float64 `json:"potential_saving"`
	PotentialSavingPct float64 `json:"potential_saving_pct"`
}

// SavingSimulation provides a simple calculation placeholder
func (h *Handler) SavingSimulation(c *gin.Context) {
	ctx := c.Request.Context()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("Panic occurred: %v", r))
			c.JSON(http.StatusInternalServerError, res)
		}
	}()

	var req savingSimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := transformer.ExceptionResponse(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Placeholder values; replace with real computation using pricing tiers and seat counts
	current := 10_000_000.0
	licenses := 2_000_000.0
	switchCost := 500_000.0
	proposed := licenses + switchCost
	saving := current - proposed
	pct := 0.0
	if current > 0 {
		pct = saving / current
	}

	output := transformer.SuccessResponse(http.StatusOK, savingSimulationResponse{
		CurrentCostYear:    current,
		ProposedLicenses:   licenses,
		SwitchingCost:      switchCost,
		ProposedTotal:      proposed,
		PotentialSaving:    saving,
		PotentialSavingPct: pct,
	})
	c.JSON(http.StatusOK, output)
}
