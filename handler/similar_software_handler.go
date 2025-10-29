package handler

import (
	dbRepo "asset-pulse-api/repositories/database"
	"asset-pulse-api/services/ai"
	"asset-pulse-api/utils/logger"
	"asset-pulse-api/utils/transformer"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// SimilarSoftwareClusterResponse represents a cluster of similar software
type SimilarSoftwareClusterResponse struct {
	Key             string                       `json:"key"`               // feature_cluster_key
	Name            string                       `json:"name"`              // Group name from AI
	Description     string                       `json:"description"`       // Description from AI
	CommonFeatures  []string                     `json:"common_features"`   // Common features from AI
	Apps            []SimilarSoftwareAppResponse `json:"apps"`              // Apps in this cluster
	TotalUsers      int                          `json:"total_users"`       // Sum of all users
	CurrentCostYear float64                      `json:"current_cost_year"` // Current total cost per year
	PotentialSaving *float64                     `json:"potential_saving"`  // Potential saving if consolidated
	Subsidiaries    []string                     `json:"subsidiaries"`      // All subsidiaries using apps in this cluster
}

// SimilarSoftwareAppResponse represents an app within a cluster
type SimilarSoftwareAppResponse struct {
	AppID           int64                              `json:"app_id"`
	Name            string                             `json:"name"`
	VendorID        *int64                             `json:"vendor_id"`
	VendorName      *string                            `json:"vendor_name"`
	Category        *string                            `json:"category"`
	Users           int                                `json:"users"`          // Total users across all subsidiaries
	PricePerSeat    float64                            `json:"price_per_seat"` // THB per seat per period
	BillingPeriod   string                             `json:"billing_period"` // monthly/yearly
	Currency        string                             `json:"currency"`
	CostYear        float64                            `json:"cost_year"`                  // Total cost per year
	Subsidiaries    []string                           `json:"subsidiaries"`               // Subsidiaries using this app
	Details         []SimilarSoftwareAppDetailResponse `json:"details"`                    // Per-subsidiary details
	SimilarityScore *float64                           `json:"similarity_score,omitempty"` // AI similarity score if from search
}

// SimilarSoftwareAppDetailResponse represents app usage per subsidiary
type SimilarSoftwareAppDetailResponse struct {
	Subsidiary    string  `json:"subsidiary"`
	Users         int     `json:"users"`
	PricePerSeat  float64 `json:"price_per_seat"`
	BillingPeriod string  `json:"billing_period"`
	CostYear      float64 `json:"cost_year"`
}

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
	selectedSubsidiaries := c.QueryArray("subsidiaries")
	appNameQuery := c.Query("app_name")
	appIDQuery := c.Query("app_id")

	var appName *string
	if appNameQuery != "" {
		appName = &appNameQuery
	}

	var appID *int64
	if appIDQuery != "" {
		if parsedAppID, err := strconv.ParseInt(appIDQuery, 10, 64); err == nil {
			appID = &parsedAppID
		}
	}

	// Step 1: Get detailed software licenses from database
	// NOTE: When app_name is provided, don't filter at DB level - let AI search filter instead
	// This allows us to get similar apps via AI service first, then filter
	var dbAppName *string
	if appID == nil {
		// Only skip DB filtering if we're not querying a specific app_id
		// We'll use AI search results to filter instead
		dbAppName = nil
	} else {
		// If app_id is provided, keep app_name filter for datatable use case
		dbAppName = appName
	}
	licenses, err := h.dbRepo.GetSimilarSoftwareLicenses(ctx, companyCode, dbAppName, selectedSubsidiaries, appID)
	if err != nil {
		res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
		logger.Error(ctx, fmt.Sprintf("Get similar software licenses error: %v", err))
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Step 2: If app_name is provided, use catalog search to find similar apps and filter licenses
	var searchResults []ai.SearchResult
	var similarAppNames map[string]bool    // Set of app names from search results
	var similarityThreshold float64 = 50.0 // Lower threshold to include more similar apps

	if appName != nil && *appName != "" {
		if h.catalogSearchService != nil {
			searchReq := &ai.CatalogSearchRequest{
				Query: *appName,
				TopK:  50, // Get more results for better coverage
			}
			searchResp, searchErr := h.catalogSearchService.Search(ctx, searchReq)
			if searchErr == nil && searchResp != nil {
				searchResults = searchResp.Results
				// Build set of similar app names (using InputName which is the original DB name)
				// Include all search results to get more apps, then group them
				similarAppNames = make(map[string]bool)
				for _, result := range searchResults {
					// Lower threshold to include more similar apps
					if result.Similarity >= similarityThreshold {
						similarAppNames[result.InputName] = true
						// Also try matching by normalized product name if available
						if result.ProductName != "" && result.ProductName != result.InputName {
							similarAppNames[result.ProductName] = true
						}
						// Include vendor name as fallback matching
						if result.Vendor != "" {
							// Try matching by vendor + product name pattern
							vendorProduct := result.Vendor + " " + result.ProductName
							similarAppNames[vendorProduct] = true
						}
					}
				}
				logger.Info(ctx, fmt.Sprintf("AI catalog search: Found %d similar apps (similarity >= %.0f%%) from %d search results", len(similarAppNames), similarityThreshold, len(searchResults)))
			} else if searchErr != nil {
				logger.Error(ctx, fmt.Sprintf("Catalog search error: %v", searchErr))
			}
		}
	}

	// Step 3: Filter licenses by similar app names if search was performed
	var filteredLicenses []dbRepo.SimilarSoftwareLicense
	if len(similarAppNames) > 0 {
		// Create multiple lookup maps for better matching
		similarAppNamesLower := make(map[string]bool)
		similarAppNamesWords := make(map[string]bool) // Individual words for better matching
		for name := range similarAppNames {
			similarAppNamesLower[strings.ToLower(name)] = true
			// Split into words for token matching
			words := strings.Fields(strings.ToLower(name))
			for _, word := range words {
				if len(word) > 2 { // Ignore very short words
					similarAppNamesWords[word] = true
				}
			}
		}

		// Filter licenses to include only apps found in search results
		// Use multiple matching strategies
		for _, lic := range licenses {
			appNameLower := strings.ToLower(lic.AppName)
			matched := false

			// Strategy 1: Exact match (case-insensitive)
			if similarAppNames[lic.AppName] || similarAppNamesLower[appNameLower] {
				matched = true
			} else {
				// Strategy 2: Partial substring matching
				for searchName := range similarAppNames {
					searchNameLower := strings.ToLower(searchName)
					if strings.Contains(appNameLower, searchNameLower) ||
						strings.Contains(searchNameLower, appNameLower) {
						matched = true
						break
					}
				}

				// Strategy 3: Word token matching (if no match yet)
				if !matched {
					licWords := strings.Fields(appNameLower)
					for _, licWord := range licWords {
						if len(licWord) > 2 && similarAppNamesWords[licWord] {
							matched = true
							break
						}
					}
				}
			}

			if matched {
				filteredLicenses = append(filteredLicenses, lic)
			}
		}
		logger.Info(ctx, fmt.Sprintf("AI similarity search: Found %d similar app names from catalog, filtered to %d licenses from %d total in database", len(similarAppNames), len(filteredLicenses), len(licenses)))
	} else if appName != nil && *appName != "" {
		// If search was attempted but no results, log and return empty
		logger.Info(ctx, fmt.Sprintf("AI similarity search returned no results for: %s", *appName))
		if len(licenses) > 0 {
			// Fallback: try simple name matching if AI search failed
			searchLower := strings.ToLower(*appName)
			for _, lic := range licenses {
				if strings.Contains(strings.ToLower(lic.AppName), searchLower) {
					filteredLicenses = append(filteredLicenses, lic)
				}
			}
			logger.Info(ctx, fmt.Sprintf("Fallback: Matched %d licenses using simple name matching", len(filteredLicenses)))
		}
	} else {
		// No search performed, use all licenses
		filteredLicenses = licenses
	}

	if len(filteredLicenses) == 0 {
		// If subsidiaries were requested but no data found, return 500 error
		if len(selectedSubsidiaries) > 0 {
			err := fmt.Errorf("no data found for requested subsidiaries: %v", selectedSubsidiaries)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("No data found for requested subsidiaries: %v", selectedSubsidiaries))
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		output := transformer.SuccessResponse(http.StatusOK, []SimilarSoftwareClusterResponse{})
		c.JSON(http.StatusOK, output)
		return
	}

	// Step 4: Build category map from catalog search results
	// If we have search results, use their categories for grouping
	categoryMap := make(map[string]string) // app name -> category
	for _, result := range searchResults {
		if result.Category != "" {
			categoryMap[result.InputName] = result.Category
			// Also map by product name if different
			if result.ProductName != "" && result.ProductName != result.InputName {
				categoryMap[result.ProductName] = result.Category
			}
		}
	}

	// Step 5: If app_id is provided, return only that app's data (for datatable)
	// Otherwise, return grouped clusters
	if appID != nil {
		logger.Info(ctx, fmt.Sprintf("Processing app_id=%d with %d licenses, requested subsidiaries: %v", *appID, len(licenses), selectedSubsidiaries))
		// Group licenses by app_id (should be only one app since we filtered)
		appMap := make(map[int64]*SimilarSoftwareAppResponse)
		for _, lic := range licenses {
			app, exists := appMap[lic.AppID]
			if !exists {
				app = &SimilarSoftwareAppResponse{
					AppID:         lic.AppID,
					Name:          lic.AppName,
					VendorID:      lic.VendorID,
					VendorName:    lic.VendorName,
					Category:      lic.AppCategory,
					Users:         0,
					PricePerSeat:  lic.PricePerSeat,
					BillingPeriod: lic.BillingPeriod,
					Currency:      lic.Currency,
					CostYear:      0,
					Subsidiaries:  []string{},
					Details:       []SimilarSoftwareAppDetailResponse{},
				}
				appMap[lic.AppID] = app
			}

			// Calculate cost per year for this subsidiary
			costYear := lic.PricePerSeat * float64(lic.Users)
			if lic.BillingPeriod == "monthly" {
				costYear = costYear * 12
			}

			// Add subsidiary detail
			detail := SimilarSoftwareAppDetailResponse{
				Subsidiary:    lic.CompanyCode,
				Users:         lic.Users,
				PricePerSeat:  lic.PricePerSeat,
				BillingPeriod: lic.BillingPeriod,
				CostYear:      costYear,
			}
			app.Details = append(app.Details, detail)
			app.Users += lic.Users
			app.CostYear += costYear

			// Add subsidiary to app's list
			if !contains(app.Subsidiaries, lic.CompanyCode) {
				app.Subsidiaries = append(app.Subsidiaries, lic.CompanyCode)
			}
		}

		// Convert map to slice
		appList := make([]SimilarSoftwareAppResponse, 0, len(appMap))
		for _, app := range appMap {
			logger.Info(ctx, fmt.Sprintf("App %s (%d) has %d details for subsidiaries: %v", app.Name, app.AppID, len(app.Details), app.Subsidiaries))
			appList = append(appList, *app)
		}

		// Return as cluster format for compatibility
		if len(appList) > 0 {
			cluster := SimilarSoftwareClusterResponse{
				Key:             "selected_app",
				Name:            appList[0].Name,
				Description:     fmt.Sprintf("Details for %s", appList[0].Name),
				CommonFeatures:  []string{},
				Apps:            appList,
				TotalUsers:      appList[0].Users,
				CurrentCostYear: appList[0].CostYear,
				Subsidiaries:    appList[0].Subsidiaries,
			}

			// Check if all requested subsidiaries are present in the result
			if len(selectedSubsidiaries) > 0 {
				resultSubsidiaries := make(map[string]bool)
				for _, sub := range appList[0].Subsidiaries {
					resultSubsidiaries[sub] = true
				}

				// Check which requested subsidiaries are missing from the result
				missingSubsidiaries := []string{}
				for _, requestedSub := range selectedSubsidiaries {
					if !resultSubsidiaries[requestedSub] {
						missingSubsidiaries = append(missingSubsidiaries, requestedSub)
					}
				}

				// If any requested subsidiaries are missing, return 500 error
				if len(missingSubsidiaries) > 0 {
					err := fmt.Errorf("no data found for requested subsidiaries: %v (found: %v)", missingSubsidiaries, appList[0].Subsidiaries)
					res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
					logger.Error(ctx, fmt.Sprintf("No data found for app_id=%d with requested subsidiaries: %v (found: %v)", *appID, missingSubsidiaries, appList[0].Subsidiaries))
					c.JSON(http.StatusInternalServerError, res)
					return
				}
			}

			output := transformer.SuccessResponse(http.StatusOK, []SimilarSoftwareClusterResponse{cluster})
			c.JSON(http.StatusOK, output)
			return
		}

		// Empty response if no data
		// If subsidiaries were requested but no data found, return 500 error
		if len(selectedSubsidiaries) > 0 {
			err := fmt.Errorf("no data found for app_id=%d with requested subsidiaries: %v", *appID, selectedSubsidiaries)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("No data found for app_id=%d with requested subsidiaries: %v", *appID, selectedSubsidiaries))
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		output := transformer.SuccessResponse(http.StatusOK, []SimilarSoftwareClusterResponse{})
		c.JSON(http.StatusOK, output)
		return
	}

	// Step 5: Transform licenses to clusters using category-based grouping
	clusters := h.transformToClustersByCategory(filteredLicenses, searchResults, categoryMap)

	// Check if all requested subsidiaries are present in the cluster results
	if len(selectedSubsidiaries) > 0 {
		// Collect all subsidiaries from all clusters
		resultSubsidiaries := make(map[string]bool)
		for _, cluster := range clusters {
			for _, sub := range cluster.Subsidiaries {
				resultSubsidiaries[sub] = true
			}
		}

		// Check which requested subsidiaries are missing from the result
		missingSubsidiaries := []string{}
		for _, requestedSub := range selectedSubsidiaries {
			if !resultSubsidiaries[requestedSub] {
				missingSubsidiaries = append(missingSubsidiaries, requestedSub)
			}
		}

		// If any requested subsidiaries are missing, return 500 error
		if len(missingSubsidiaries) > 0 {
			// Get list of found subsidiaries for logging
			foundSubsidiaries := make([]string, 0, len(resultSubsidiaries))
			for sub := range resultSubsidiaries {
				foundSubsidiaries = append(foundSubsidiaries, sub)
			}
			err := fmt.Errorf("no data found for requested subsidiaries: %v (found: %v)", missingSubsidiaries, foundSubsidiaries)
			res := transformer.ExceptionResponse(http.StatusInternalServerError, err)
			logger.Error(ctx, fmt.Sprintf("No data found for requested subsidiaries: %v (found: %v)", missingSubsidiaries, foundSubsidiaries))
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	output := transformer.SuccessResponse(http.StatusOK, clusters)
	c.JSON(http.StatusOK, output)
}

// transformToClustersByCategory groups licenses by category from catalog search or database
func (h *Handler) transformToClustersByCategory(
	licenses []dbRepo.SimilarSoftwareLicense,
	searchResults []ai.SearchResult,
	categoryMap map[string]string,
) []SimilarSoftwareClusterResponse {
	clustersMap := make(map[string]*SimilarSoftwareClusterResponse)

	// Build similarity score map from search results (for apps found via AI search)
	similarityMap := make(map[string]float64) // key: app name -> similarity score
	for _, result := range searchResults {
		similarityMap[result.InputName] = result.Similarity
		// Also map by product name if different
		if result.ProductName != "" && result.ProductName != result.InputName {
			similarityMap[result.ProductName] = result.Similarity
		}
	}

	// Group licenses by category
	for _, lic := range licenses {
		// Get category from catalog search map, or DB category, or fallback to "General"
		category := "General"

		// Priority 1: Category from catalog search results
		if cat, found := categoryMap[lic.AppName]; found && cat != "" {
			category = cat
		} else if lic.AppCategory != nil && *lic.AppCategory != "" {
			// Priority 2: Category from database
			category = *lic.AppCategory
		}

		clusterKey := strings.ToLower(strings.ReplaceAll(category, " ", "_"))

		cluster, exists := clustersMap[clusterKey]
		if !exists {
			// Create cluster with meaningful description based on category
			description := fmt.Sprintf("Similar %s applications", category)
			commonFeatures := []string{"Core functionality", "Enterprise features", "Team collaboration"}

			// Add category-specific features
			switch category {
			case "Design Tools":
				commonFeatures = []string{"Visual design", "Prototyping", "Collaboration", "Design systems"}
			case "Developer Tools":
				commonFeatures = []string{"Code management", "Version control", "CI/CD", "Development workflow"}
			case "Project Management":
				commonFeatures = []string{"Task tracking", "Agile boards", "Sprint management", "Team collaboration"}
			case "Business Intelligence":
				commonFeatures = []string{"Data visualization", "Reporting", "Analytics", "Dashboards"}
			case "Communication & Collaboration":
				commonFeatures = []string{"Messaging", "Video conferencing", "File sharing", "Team collaboration"}
			case "IT Service Management":
				commonFeatures = []string{"Incident management", "Service desk", "Ticketing", "IT operations"}
			case "Cloud Platforms":
				commonFeatures = []string{"Infrastructure", "Deployment", "Scalability", "Cloud services"}
			case "Productivity Suites":
				commonFeatures = []string{"Office applications", "Document management", "Email", "Calendar"}
			}

			cluster = &SimilarSoftwareClusterResponse{
				Key:            clusterKey,
				Name:           category,
				Description:    description,
				CommonFeatures: commonFeatures,
				Apps:           []SimilarSoftwareAppResponse{},
				Subsidiaries:   []string{},
			}
			clustersMap[clusterKey] = cluster
		}

		h.addLicenseToCluster(cluster, lic, similarityMap)
	}

	// Calculate metrics for all clusters
	for _, cluster := range clustersMap {
		h.calculateClusterMetrics(cluster)
	}

	// Convert map to slice
	result := make([]SimilarSoftwareClusterResponse, 0, len(clustersMap))
	for _, cluster := range clustersMap {
		result = append(result, *cluster)
	}

	return result
}

// addLicenseToCluster adds a license to the appropriate app in a cluster
func (h *Handler) addLicenseToCluster(cluster *SimilarSoftwareClusterResponse, lic dbRepo.SimilarSoftwareLicense, similarityMap map[string]float64) {
	// Find or create app in cluster
	var app *SimilarSoftwareAppResponse
	appIdx := -1
	for i := range cluster.Apps {
		if cluster.Apps[i].AppID == lic.AppID {
			app = &cluster.Apps[i]
			appIdx = i
			break
		}
	}

	if app == nil {
		// Create new app
		app = &SimilarSoftwareAppResponse{
			AppID:         lic.AppID,
			Name:          lic.AppName,
			VendorID:      lic.VendorID,
			VendorName:    lic.VendorName,
			Category:      lic.AppCategory,
			Users:         0,
			PricePerSeat:  lic.PricePerSeat,
			BillingPeriod: lic.BillingPeriod,
			Currency:      lic.Currency,
			CostYear:      0,
			Subsidiaries:  []string{},
			Details:       []SimilarSoftwareAppDetailResponse{},
		}
		// Set similarity score if available from AI search
		if similarityMap != nil {
			if score, found := similarityMap[lic.AppName]; found {
				app.SimilarityScore = &score
			}
		}
		cluster.Apps = append(cluster.Apps, *app)
		appIdx = len(cluster.Apps) - 1
	} else if similarityMap != nil {
		// Update similarity score if app already exists and we have score for it
		if score, found := similarityMap[lic.AppName]; found {
			cluster.Apps[appIdx].SimilarityScore = &score
		}
	}

	// Calculate cost per year for this subsidiary
	costYear := lic.PricePerSeat * float64(lic.Users)
	if lic.BillingPeriod == "monthly" {
		costYear = costYear * 12
	}

	// Add subsidiary detail
	detail := SimilarSoftwareAppDetailResponse{
		Subsidiary:    lic.CompanyCode,
		Users:         lic.Users,
		PricePerSeat:  lic.PricePerSeat,
		BillingPeriod: lic.BillingPeriod,
		CostYear:      costYear,
	}
	cluster.Apps[appIdx].Details = append(cluster.Apps[appIdx].Details, detail)
	cluster.Apps[appIdx].Users += lic.Users
	cluster.Apps[appIdx].CostYear += costYear

	// Add subsidiary to app's list
	if !contains(cluster.Apps[appIdx].Subsidiaries, lic.CompanyCode) {
		cluster.Apps[appIdx].Subsidiaries = append(cluster.Apps[appIdx].Subsidiaries, lic.CompanyCode)
	}

	// Add subsidiary to cluster's list
	if !contains(cluster.Subsidiaries, lic.CompanyCode) {
		cluster.Subsidiaries = append(cluster.Subsidiaries, lic.CompanyCode)
	}
}

// calculateClusterMetrics calculates total metrics for a cluster
func (h *Handler) calculateClusterMetrics(cluster *SimilarSoftwareClusterResponse) {
	totalUsers := 0
	totalCost := 0.0

	for i := range cluster.Apps {
		totalUsers += cluster.Apps[i].Users
		totalCost += cluster.Apps[i].CostYear
	}

	cluster.TotalUsers = totalUsers
	cluster.CurrentCostYear = totalCost
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
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
