package cmd

import (
	"fmt"

	"github.com/roboalchemist/exa-cli/pkg/api"
	"github.com/roboalchemist/exa-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	similarNumResults    int
	similarExcludeSource bool
	similarIncDomains    []string
	similarExcDomains    []string
	similarStartDate     string
	similarEndDate       string
	similarText          bool
	similarHighlights    bool
	similarCategory      string
)

var similarCmd = &cobra.Command{
	Use:   "similar [url]",
	Short: "Find pages similar to a URL",
	Long: `Find web pages similar to the given URL.

Uses Exa's neural search to find semantically similar content.

Examples:
  exa similar "https://arxiv.org/abs/2307.06435"
  exa similar "https://example.com" -n 20
  exa similar "https://blog.example.com" --exclude-source
  exa similar "https://example.com" --include-domains arxiv.org,scholar.google.com
  exa similar "https://example.com" --json`,
	Args: cobra.ExactArgs(1),
	RunE: runSimilar,
}

func init() {
	f := similarCmd.Flags()
	f.IntVarP(&similarNumResults, "num-results", "n", 10, "Max results")
	f.BoolVar(&similarExcludeSource, "exclude-source", false, "Exclude the source domain from results")
	f.StringSliceVar(&similarIncDomains, "include-domains", nil, "Only include these domains")
	f.StringSliceVar(&similarExcDomains, "exclude-domains", nil, "Exclude these domains")
	f.StringVar(&similarStartDate, "start-date", "", "Published after (YYYY-MM-DD)")
	f.StringVar(&similarEndDate, "end-date", "", "Published before (YYYY-MM-DD)")
	f.BoolVar(&similarText, "text", false, "Include full text")
	f.BoolVar(&similarHighlights, "highlights", false, "Include highlights")
	f.StringVar(&similarCategory, "category", "", "Category filter")

	rootCmd.AddCommand(similarCmd)
}

func runSimilar(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	req := &api.FindSimilarRequest{
		URL:                 args[0],
		NumResults:          similarNumResults,
		ExcludeSourceDomain: similarExcludeSource,
	}

	if len(similarIncDomains) > 0 {
		req.IncludeDomains = similarIncDomains
	}
	if len(similarExcDomains) > 0 {
		req.ExcludeDomains = similarExcDomains
	}
	if similarStartDate != "" {
		req.StartPublishedDate = similarStartDate + "T00:00:00.000Z"
	}
	if similarEndDate != "" {
		req.EndPublishedDate = similarEndDate + "T00:00:00.000Z"
	}
	if similarCategory != "" {
		req.Category = similarCategory
	}

	if similarText || similarHighlights {
		contents := &api.ContentsSpec{}
		if similarText {
			contents.Text = &api.TextSpec{MaxCharacters: 10000}
		}
		if similarHighlights {
			contents.Highlights = &api.HighlightsSpec{}
		}
		req.Contents = contents
	}

	resp, err := client.FindSimilar(newContext(), req)
	if err != nil {
		return err
	}

	opts := GetOutputOptions()

	if opts.Mode == output.ModeJSON {
		return output.RenderJSON(resp, opts)
	}

	td := output.TableData{
		Headers: []string{"TITLE", "URL", "DATE", "SCORE"},
	}

	for _, r := range resp.Results {
		date := ""
		if r.PublishedDate != "" && len(r.PublishedDate) >= 10 {
			date = r.PublishedDate[:10]
		}
		td.Rows = append(td.Rows, []string{truncateStr(r.Title, 50), r.URL, date, fmt.Sprintf("%.2f", r.Score)})
	}

	footer := fmt.Sprintf("%d similar pages", len(resp.Results))
	if resp.CostDollars != nil {
		footer = fmt.Sprintf("Cost: $%.4f | %s", resp.CostDollars.Total, footer)
	}
	td.Footer = footer

	return output.RenderTable(td, resp, opts)
}
