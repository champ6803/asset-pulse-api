package ai

import (
	"asset-pulse-api/entities"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// CatalogSearchService handles AI-powered catalog search
type CatalogSearchService struct {
	client      *openai.Client
	db          *gorm.DB
	catalog     []CatalogItem
	embeddings  [][]float32
	initialized bool
}

// CatalogItem represents a software catalog item
type CatalogItem struct {
	ID              int64
	InputName       string
	Vendor          string
	ProductName     string
	Functionalities string
	Embedding       []float32
}

// SearchResult represents a catalog search result
type SearchResult struct {
	ProductName     string  `json:"product_name"`
	Vendor          string  `json:"vendor"`
	Functionalities string  `json:"functionalities"`
	InputName       string  `json:"input_name"`
	Similarity      float64 `json:"similarity"`
}

// CatalogSearchRequest for searching the catalog
type CatalogSearchRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"top_k"`
}

// CatalogSearchResponse contains search results
type CatalogSearchResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
	TopK    int            `json:"top_k"`
}

const (
	embedModel     = "text-embedding-3-large"
	topKDefault    = 5
	normalizeModel = "gpt-4o-mini"
)

// NewCatalogSearchService creates a new catalog search service
func NewCatalogSearchService(apiKey string, db *gorm.DB) *CatalogSearchService {
	client := openai.NewClient(apiKey)
	return &CatalogSearchService{
		client:      client,
		db:          db,
		catalog:     []CatalogItem{},
		embeddings:  [][]float32{},
		initialized: false,
	}
}

// Initialize loads and processes the catalog from database
func (s *CatalogSearchService) Initialize(ctx context.Context) error {
	if s.initialized {
		return nil
	}

	// Fetch all apps from database
	var apps []entities.App
	if err := s.db.WithContext(ctx).Find(&apps).Error; err != nil {
		return fmt.Errorf("failed to fetch apps: %w", err)
	}

	if len(apps) == 0 {
		return fmt.Errorf("no apps found in database")
	}

	fmt.Printf("Loading %d apps from database...\n", len(apps))

	// Process each app: normalize and embed
	for i, app := range apps {
		fmt.Printf("Processing app %d/%d: %s\n", i+1, len(apps), app.Name)

		// Normalize the app description
		normalized, err := s.normalizeDescription(ctx, app.Name)
		if err != nil {
			fmt.Printf("Error normalizing %s: %v\n", app.Name, err)
			continue
		}

		// Create embedding for the normalized text
		textToEmbed := normalized.Vendor + " " + normalized.ProductName + " " + normalized.Functionalities
		embedding, err := s.createEmbedding(ctx, textToEmbed)
		if err != nil {
			fmt.Printf("Error creating embedding for %s: %v\n", app.Name, err)
			continue
		}

		item := CatalogItem{
			ID:              app.ID,
			InputName:       app.Name,
			Vendor:          normalized.Vendor,
			ProductName:     normalized.ProductName,
			Functionalities: normalized.Functionalities,
			Embedding:       embedding,
		}

		s.catalog = append(s.catalog, item)
		s.embeddings = append(s.embeddings, embedding)
	}

	fmt.Printf("Successfully loaded %d items into catalog\n", len(s.catalog))
	s.initialized = true
	return nil
}

// NormalizedOutput represents the normalized product information
type NormalizedOutput struct {
	Vendor          string `json:"vendor"`
	ProductName     string `json:"product_name"`
	Functionalities string `json:"functionalities"`
}

// normalizeDescription normalizes a product description using OpenAI
func (s *CatalogSearchService) normalizeDescription(ctx context.Context, text string) (*NormalizedOutput, error) {
	prompt := fmt.Sprintf(`Extract, normalize and summarize the following software product or license string into json text of
1. vendor - precise vendor name
2. product_name - concise, normalized commercial product name
3. functionalities - return concise, functionality list of the software. try to not mention vendor-specific wording or specific vendor tool, while ok to return wording of industry-standard protocols.

Example: 
  Input: 'Figma Organization license for PointX project 2025 3 seats purchaseID AA22343'
  Output: {"vendor": "Figma", "product_name": "Figma Design", "functionalities": "Vector drawing and editing, collaborative design, prototyping, developer handoff"}

Input: %s
Output:`, text)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       normalizeModel,
		Messages:    []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}},
		Temperature: 0.3,
	})

	if err != nil {
		return nil, fmt.Errorf("openai api error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from openai")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Parse JSON response
	var normalized NormalizedOutput
	if err := json.Unmarshal([]byte(content), &normalized); err != nil {
		// If JSON parsing fails, try to extract from markdown code blocks
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)

		if err := json.Unmarshal([]byte(content), &normalized); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	return &normalized, nil
}

// createEmbedding creates an embedding for the given text
func (s *CatalogSearchService) createEmbedding(ctx context.Context, text string) ([]float32, error) {
	resp, err := s.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(embedModel),
		Input: []string{text},
	})

	if err != nil {
		return nil, fmt.Errorf("openai embedding error: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	return resp.Data[0].Embedding, nil
}

// Search searches the catalog for similar items
func (s *CatalogSearchService) Search(ctx context.Context, req *CatalogSearchRequest) (*CatalogSearchResponse, error) {
	if !s.initialized {
		return nil, fmt.Errorf("service not initialized; call Initialize() first")
	}

	topK := req.TopK
	if topK <= 0 {
		topK = topKDefault
	}
	if topK > len(s.catalog) {
		topK = len(s.catalog)
	}

	// Normalize the input query
	normalized, err := s.normalizeDescription(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize query: %w", err)
	}

	// Create embedding for the query
	textToEmbed := normalized.Vendor + " " + normalized.ProductName + " " + normalized.Functionalities
	queryEmbedding, err := s.createEmbedding(ctx, textToEmbed)
	if err != nil {
		return nil, fmt.Errorf("failed to create query embedding: %w", err)
	}

	// Normalize query embedding
	queryNorm := s.normalizeVector(queryEmbedding)

	// Calculate similarity scores
	type scoredItem struct {
		item  CatalogItem
		score float64
		index int
	}

	scores := make([]scoredItem, len(s.catalog))
	for i, item := range s.catalog {
		// Normalize catalog item embedding
		catalogNorm := s.normalizeVector(item.Embedding)
		score := s.dotProduct(queryNorm, catalogNorm)
		scores[i] = scoredItem{
			item:  item,
			score: score,
			index: i,
		}
	}

	// Sort by score descending and take top K
	for i := 0; i < len(scores)-1; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[i].score < scores[j].score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// Convert to SearchResults
	results := make([]SearchResult, topK)
	for i := 0; i < topK; i++ {
		results[i] = SearchResult{
			ProductName:     scores[i].item.ProductName,
			Vendor:          scores[i].item.Vendor,
			Functionalities: scores[i].item.Functionalities,
			InputName:       scores[i].item.InputName,
			Similarity:      math.Round(scores[i].score*10000) / 100, // Round to 2 decimal places
		}
	}

	return &CatalogSearchResponse{
		Query:   req.Query,
		Results: results,
		TopK:    topK,
	}, nil
}

// normalizeVector normalizes a vector using L2 normalization
func (s *CatalogSearchService) normalizeVector(vec []float32) []float32 {
	sum := float32(0)
	for _, v := range vec {
		sum += v * v
	}
	norm := float32(math.Sqrt(float64(sum)))
	if norm == 0 {
		return vec
	}

	normalized := make([]float32, len(vec))
	for i, v := range vec {
		normalized[i] = v / norm
	}
	return normalized
}

// dotProduct calculates the dot product of two vectors
func (s *CatalogSearchService) dotProduct(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	sum := float32(0)
	for i := range a {
		sum += a[i] * b[i]
	}
	return float64(sum)
}
