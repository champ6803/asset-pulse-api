package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type SoftwareLicense struct {
	Name              string  `json:"name"`
	THBPricePerYear   float64 `json:"thbPricePerYear"`
	UsedByCompanyName string  `json:"usedByCompanyName"`
}

type GroupedSoftwareResponse struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Items          []map[string]any `json:"items"`
	CommonFeatures []string         `json:"commonFeatures"`
}

type SoftwareGroupingService struct {
	client       *azopenai.Client
	deploymentID string
}

func NewSoftwareGroupingService(endpoint, deploymentID string) (*SoftwareGroupingService, error) {
	// cred, err := azidentity.NewDefaultAzureCredential(nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get Azure credentials: %w", err)
	// }

	// client, err := azopenai.NewClient(endpoint, cred, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create Azure OpenAI client: %w", err)
	// }

	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	if apiKey == "" {
		panic("missing AZURE_OPENAI_API_KEY")
	}

	// create the credential
	cred := azcore.NewKeyCredential(apiKey)
	client, err := azopenai.NewClientWithKeyCredential(endpoint, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure OpenAI client: %w", err)
	}

	return &SoftwareGroupingService{
		client:       client,
		deploymentID: deploymentID,
	}, nil
}

func (s *SoftwareGroupingService) GroupSoftwareLicenses(ctx context.Context, licenses []SoftwareLicense) ([]GroupedSoftwareResponse, error) {
	inputJSON, err := json.MarshalIndent(licenses, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	systemPrompt := `
You are an expert enterprise software analyst.

You will receive a JSON array of software licenses purchased by various subsidiaries. Each license has:
- name: string, the software name (may contain typos)
- thbPricePerYear: number, the license price in Thai Baht per year
- usedByCompanyName: string, the subsidiary that purchased it

Your tasks:
1. Normalize software names (correct typos and unify variants, e.g., "Figam" -> "Figma").
2. Group the software into logical categories (e.g., "Design Tools", "Developer Tools", "Project Management", "Business Intelligence", etc.).
3. For each group, output a JSON object with the following structure:

{
  "name": string,              // Group name
  "description": string,       // Description of the group
  "items": [                   // Array of software items
    {
      "id": string,            // Unique ID for the software (can be sequential)
      "name": string,          // Normalized software name
      "image": string,         // URL to a representative image
      "licensePricePerYear": number
    }
  ],
  "commonFeatures": [string]   // List of features common to software in this group
}

Important rules:
- The output **must be valid JSON only**. Do not include explanations, comments, or extra text.
- Prefer a **top-level array of group objects**.
- Ensure each software is included exactly once under its appropriate category.
- Try to deduplicate software entries purchased by different subsidiaries.
- Correct obvious typos in software names.
- Use sequential string IDs for each software item (e.g., "1", "2", "3"...).
- Provide reasonable group descriptions and common features based on the software in the group.

Example of correct output format:

[
  {
    "name": "Design Tools",
    "description": "Software used for graphic design, prototyping, and UI/UX design.",
    "items": [
      {
        "id": "1",
        "name": "Figma",
        "image": "https://example.com/images/figma.png",
        "licensePricePerYear": 4800
      }
    ],
    "commonFeatures": ["Collaboration", "Prototyping", "Design Systems"]
  },
  {
    "name": "Project Management",
    "description": "Tools for planning, tracking, and managing projects and tasks.",
    "items": [
      {
        "id": "2",
        "name": "Jira Software",
        "image": "https://example.com/images/jira.png",
        "licensePricePerYear": 12000
      }
    ],
    "commonFeatures": ["Task Management", "Agile Boards", "Reporting"]
  }
]
`

	userPrompt := fmt.Sprintf(
		"Here is the input JSON array of software licenses:\n\n%s\n\nPlease generate the grouped structured output as specified.",
		string(inputJSON),
	)

	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{
			Content: azopenai.NewChatRequestSystemMessageContent(systemPrompt),
		},
		&azopenai.ChatRequestUserMessage{
			Content: azopenai.NewChatRequestUserMessageContent(userPrompt),
		},
	}

	req := azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &s.deploymentID,
		Temperature:    toPtr[float32](0.3),
	}

	resp, err := s.client.GetChatCompletions(ctx, req, nil)
	if err != nil {
		return nil, fmt.Errorf("azure openai api error: %w", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == nil {
		return nil, fmt.Errorf("no valid response from Azure OpenAI")
	}

	content := *resp.Choices[0].Message.Content

	var result []GroupedSoftwareResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	return result, nil
}

func toPtr[T any](v T) *T { return &v }
