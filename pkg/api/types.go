package api

import "encoding/json"

// SearchRequest is the request body for POST /search
type SearchRequest struct {
	Query              string         `json:"query"`
	Type               string         `json:"type,omitempty"`
	NumResults         int            `json:"numResults,omitempty"`
	Category           string         `json:"category,omitempty"`
	IncludeDomains     []string       `json:"includeDomains,omitempty"`
	ExcludeDomains     []string       `json:"excludeDomains,omitempty"`
	StartPublishedDate string         `json:"startPublishedDate,omitempty"`
	EndPublishedDate   string         `json:"endPublishedDate,omitempty"`
	StartCrawlDate     string         `json:"startCrawlDate,omitempty"`
	EndCrawlDate       string         `json:"endCrawlDate,omitempty"`
	IncludeText        string         `json:"includeText,omitempty"`
	ExcludeText        string         `json:"excludeText,omitempty"`
	Moderation         *bool          `json:"moderation,omitempty"`
	Contents           *ContentsSpec  `json:"contents,omitempty"`
}

// ContentsSpec specifies what content to include in results.
type ContentsSpec struct {
	Text         *TextSpec       `json:"text,omitempty"`
	Highlights   *HighlightsSpec `json:"highlights,omitempty"`
	Summary      *SummarySpec    `json:"summary,omitempty"`
	Livecrawl    string          `json:"livecrawl,omitempty"`
	Subpages     int             `json:"subpages,omitempty"`
	SubpageTarget []string       `json:"subpageTarget,omitempty"`
	Extras       *ExtrasSpec     `json:"extras,omitempty"`
}

// TextSpec configures text content retrieval.
type TextSpec struct {
	MaxCharacters int  `json:"maxCharacters,omitempty"`
	IncludeHtmlTags bool `json:"includeHtmlTags,omitempty"`
}

// HighlightsSpec configures highlight extraction.
type HighlightsSpec struct {
	NumSentences      int    `json:"numSentences,omitempty"`
	HighlightsPerURL  int    `json:"highlightsPerUrl,omitempty"`
	Query             string `json:"query,omitempty"`
}

// SummarySpec configures summary generation.
type SummarySpec struct {
	Query string `json:"query,omitempty"`
}

// ExtrasSpec configures extra content extraction.
type ExtrasSpec struct {
	Links      bool `json:"links,omitempty"`
	ImageLinks bool `json:"imageLinks,omitempty"`
}

// SearchResponse is the response from POST /search
type SearchResponse struct {
	RequestID       string         `json:"requestId,omitempty"`
	ResolvedSearchType string     `json:"resolvedSearchType,omitempty"`
	Results         []SearchResult `json:"results"`
	AutopromptString string       `json:"autopromptString,omitempty"`
	CostDollars     *CostInfo      `json:"costDollars,omitempty"`
}

// SearchResult is a single search result.
type SearchResult struct {
	Title          string   `json:"title"`
	URL            string   `json:"url"`
	ID             string   `json:"id"`
	PublishedDate  string   `json:"publishedDate,omitempty"`
	Author         string   `json:"author,omitempty"`
	Score          float64  `json:"score"`
	Image          string   `json:"image,omitempty"`
	Favicon        string   `json:"favicon,omitempty"`
	Text           string   `json:"text,omitempty"`
	Highlights     []string `json:"highlights,omitempty"`
	HighlightScores []float64 `json:"highlightScores,omitempty"`
	Summary        string   `json:"summary,omitempty"`
	Subpages       []Subpage `json:"subpages,omitempty"`
}

// Subpage is a crawled subpage.
type Subpage struct {
	URL  string `json:"url"`
	Text string `json:"text,omitempty"`
}

// CostInfo contains billing information.
type CostInfo struct {
	Total     float64          `json:"total"`
	Search    *CostBreakdown   `json:"search,omitempty"`
	Contents  *CostBreakdown   `json:"contents,omitempty"`
	Neural    *CostBreakdown   `json:"neural,omitempty"`
}

// CostBreakdown is a detailed cost line item.
type CostBreakdown struct {
	Amount float64 `json:"amount"`
}

// ContentsRequest is the request body for POST /contents
type ContentsRequest struct {
	IDs      []string      `json:"ids,omitempty"`
	URLs     []string      `json:"urls,omitempty"`
	Text     *TextSpec     `json:"text,omitempty"`
	Highlights *HighlightsSpec `json:"highlights,omitempty"`
	Summary  *SummarySpec  `json:"summary,omitempty"`
	Livecrawl string       `json:"livecrawl,omitempty"`
	Subpages int           `json:"subpages,omitempty"`
}

// ContentsResponse is the response from POST /contents
type ContentsResponse struct {
	Results     []SearchResult `json:"results"`
	CostDollars *CostInfo      `json:"costDollars,omitempty"`
}

// FindSimilarRequest is the request body for POST /findSimilar
type FindSimilarRequest struct {
	URL                string   `json:"url"`
	NumResults         int      `json:"numResults,omitempty"`
	IncludeDomains     []string `json:"includeDomains,omitempty"`
	ExcludeDomains     []string `json:"excludeDomains,omitempty"`
	StartPublishedDate string   `json:"startPublishedDate,omitempty"`
	EndPublishedDate   string   `json:"endPublishedDate,omitempty"`
	ExcludeSourceDomain bool   `json:"excludeSourceDomain,omitempty"`
	Category           string   `json:"category,omitempty"`
	Contents           *ContentsSpec `json:"contents,omitempty"`
}

// FindSimilarResponse is the response from POST /findSimilar
type FindSimilarResponse struct {
	RequestID       string         `json:"requestId,omitempty"`
	Results         []SearchResult `json:"results"`
	AutopromptString string       `json:"autopromptString,omitempty"`
	CostDollars     *CostInfo      `json:"costDollars,omitempty"`
}

// AnswerRequest is the request body for POST /answer
type AnswerRequest struct {
	Query          string        `json:"query"`
	Text           bool          `json:"text,omitempty"`
	Model          string        `json:"model,omitempty"`
	OutputSchema   interface{}   `json:"outputSchema,omitempty"`
	StreamOutput   bool          `json:"stream,omitempty"`
}

// AnswerResponse is the response from POST /answer
type AnswerResponse struct {
	RequestID   string         `json:"requestId,omitempty"`
	Answer      string         `json:"answer"`
	Citations   []SearchResult `json:"citations,omitempty"`
	CostDollars *CostInfo      `json:"costDollars,omitempty"`
}

// AnswerStreamChunk represents a chunk from a streaming answer response.
type AnswerStreamChunk struct {
	Type        string         `json:"type,omitempty"`
	Text        string         `json:"text,omitempty"`
	Answer      string         `json:"answer,omitempty"`
	Citations   []SearchResult `json:"citations,omitempty"`
	CostDollars *CostInfo      `json:"costDollars,omitempty"`
}

// ContextRequest is the request body for POST /context
type ContextRequest struct {
	Query     string      `json:"query"`
	TokensNum interface{} `json:"tokensNum,omitempty"`
}

// ContextResponse is the response from POST /context
type ContextResponse struct {
	RequestID    string          `json:"requestId,omitempty"`
	Query        string          `json:"query,omitempty"`
	Context      string          `json:"response"`
	CostDollars  json.RawMessage `json:"costDollars,omitempty"`
	ResultsCount int             `json:"resultsCount,omitempty"`
	SearchTime   float64         `json:"searchTime,omitempty"`
	OutputTokens int             `json:"outputTokens,omitempty"`
	ParsedCost   *CostInfo       `json:"-"`
}

// GetCost parses the cost info, handling both string and object formats.
func (r *ContextResponse) GetCost() *CostInfo {
	if r.ParsedCost != nil {
		return r.ParsedCost
	}
	if len(r.CostDollars) == 0 {
		return nil
	}
	var cost CostInfo
	// Try direct object first
	if err := json.Unmarshal(r.CostDollars, &cost); err == nil {
		r.ParsedCost = &cost
		return r.ParsedCost
	}
	// Try string-encoded JSON
	var s string
	if err := json.Unmarshal(r.CostDollars, &s); err == nil {
		if err := json.Unmarshal([]byte(s), &cost); err == nil {
			r.ParsedCost = &cost
			return r.ParsedCost
		}
	}
	return nil
}

// UsageResponse is the response from GET /team-management/api-keys/{id}/usage
type UsageResponse struct {
	Usage []UsageEntry `json:"usage"`
}

// UsageEntry is a single usage data point.
type UsageEntry struct {
	Date          string  `json:"date"`
	RequestCount  int     `json:"requestCount"`
	CreditUsage   float64 `json:"creditUsage"`
}

// APIKeyInfo contains API key metadata.
type APIKeyInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// APIKeysResponse is the response from GET /team-management/api-keys
type APIKeysResponse struct {
	APIKeys []APIKeyInfo `json:"apiKeys"`
}
