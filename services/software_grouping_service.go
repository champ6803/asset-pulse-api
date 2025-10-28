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
	THBPricePerYear   float64 `json:"thb_price_per_year"`
	UsedByCompanyName string  `json:"used_by_company_name"`
}

type GroupedSoftwareResponse struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Items          []map[string]any `json:"items"`
	CommonFeatures []string         `json:"common_features"`
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
- name: string — the software name (may contain typos or variations)
- thb_price_per_year: number — the annual license cost in Thai Baht
- used_by_company_name: string — the subsidiary or company that purchased it

Your objectives:
1. **Normalize software names**: Correct typos and unify naming variants (e.g., "Figam" → "Figma", "Jira SW" → "Jira Software").
2. **Group software** into logical enterprise categories such as:
   - Design Tools
   - Developer Tools
   - Project Management
   - Business Intelligence
   - Communication & Collaboration
   - IT Service Management
   - Cloud Platforms
   - Productivity Suites
   (Add more if appropriate based on the data.)
3. **Preserve company usage**: If multiple subsidiaries purchased the same software, include **each usage as a separate item** within the same group, with their own "used_by_company_name" and "thb_price_per_year".
4. **Provide accurate representative images**: Use **real, verifiable logo or product image URLs** from the public internet (e.g., company official websites, Wikipedia, or reputable sources like logos-world.net, Wikimedia, or product-cdn URLs).

Output Format:
You must output a valid JSON array of group objects with this exact structure:

[
  {
    "name": string,               // Group name
    "description": string,        // Description of the group
    "items": [
      {
        "name": string,           // Normalized software name
        "image": string,          // Realistic valid image URL from the internet
        "thb_price_per_year": number,
        "used_by_company_name": string
      }
    ],
    "common_features": [string]   // List of features common to this software group
  }
]

Additional Rules:
- The **output must be valid JSON only** (no comments, explanations, or markdown).
- The **same software name** should always be normalized identically across subsidiaries.
- Each "items" array should contain **distinct entries** for each subsidiary that uses the same software.
- Do not fabricate image URLs — use **real images from credible online sources**.
- Each group must have a concise but informative **description** and **at least three common features**.

Example Output:

[
  {
    "name": "Design Tools",
    "description": "Software used for UI/UX design, prototyping, and collaboration between designers and developers.",
    "items": [
      {
        "name": "Figma",
        "image": "https://upload.wikimedia.org/wikipedia/commons/3/33/Figma-logo.svg",
        "thb_price_per_year": 4800,
        "used_by_company_name": "Company A"
      },
      {
        "name": "Figma",
        "image": "https://upload.wikimedia.org/wikipedia/commons/3/33/Figma-logo.svg",
        "thb_price_per_year": 5100,
        "used_by_company_name": "Company B"
      }
    ],
    "common_features": ["Collaboration", "Design Systems", "Prototyping"]
  },
  {
    "name": "Project Management",
    "description": "Tools used to plan, track, and manage tasks across teams.",
    "items": [
      {
        "name": "Jira Software",
        "image": "https://upload.wikimedia.org/wikipedia/commons/8/8e/Jira_Software_logo.svg",
        "thb_price_per_year": 12000,
        "used_by_company_name": "Company C"
      }
    ],
    "common_features": ["Agile Boards", "Task Tracking", "Sprint Management"]
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
