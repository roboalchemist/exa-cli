package cmd

import (
	"fmt"

	"github.com/roboalchemist/exa-cli/pkg/api"
	"github.com/roboalchemist/exa-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	contentsText       bool
	contentsTextMax    int
	contentsHighlights bool
	contentsSummary    bool
	contentsMaxAge     int
	contentsSubpages   int
)

var contentsCmd = &cobra.Command{
	Use:   "contents [urls...]",
	Short: "Get page contents by URL",
	Long: `Retrieve the contents of web pages by URL.

Returns text, highlights, and summaries for the specified URLs.

Examples:
  exa contents https://example.com
  exa contents https://example.com https://another.com
  exa contents https://example.com --highlights --summary
  exa contents https://example.com --text-max-chars 5000
  exa contents https://example.com --json`,
	Args: cobra.MinimumNArgs(1),
	RunE: runContents,
}

func init() {
	f := contentsCmd.Flags()
	f.BoolVar(&contentsText, "text", true, "Include full text")
	f.IntVar(&contentsTextMax, "text-max-chars", 10000, "Max chars for text")
	f.BoolVar(&contentsHighlights, "highlights", false, "Include highlights")
	f.BoolVar(&contentsSummary, "summary", false, "Include summary")
	f.IntVar(&contentsMaxAge, "max-age-hours", -1, "Content freshness (-1=cache, 0=always livecrawl)")
	f.IntVar(&contentsSubpages, "subpages", 0, "Subpages to crawl")

	rootCmd.AddCommand(contentsCmd)
}

func runContents(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	req := &api.ContentsRequest{
		URLs: args,
	}

	if contentsText {
		req.Text = &api.TextSpec{MaxCharacters: contentsTextMax}
	}
	if contentsHighlights {
		req.Highlights = &api.HighlightsSpec{}
	}
	if contentsSummary {
		req.Summary = &api.SummarySpec{}
	}
	if contentsMaxAge >= 0 {
		if contentsMaxAge == 0 {
			req.Livecrawl = "always"
		} else {
			req.Livecrawl = "fallback"
		}
	}
	if contentsSubpages > 0 {
		req.Subpages = contentsSubpages
	}

	resp, err := client.GetContents(newContext(), req)
	if err != nil {
		return err
	}

	opts := GetOutputOptions()

	if opts.Mode == output.ModeJSON {
		return output.RenderJSON(resp, opts)
	}

	// Table mode: show title and URL, then text below
	td := output.TableData{
		Headers: []string{"TITLE", "URL"},
	}
	for _, r := range resp.Results {
		td.Rows = append(td.Rows, []string{truncateStr(r.Title, 60), r.URL})
	}

	footer := fmt.Sprintf("%d pages", len(resp.Results))
	if resp.CostDollars != nil {
		footer = fmt.Sprintf("Cost: $%.4f | %s", resp.CostDollars.Total, footer)
	}
	td.Footer = footer

	if err := output.RenderTable(td, resp, opts); err != nil {
		return err
	}

	// Print text content below table for each result
	for _, r := range resp.Results {
		if r.Text != "" {
			fmt.Printf("\n--- %s ---\n%s\n", r.URL, truncateStr(r.Text, 2000))
		}
		if r.Summary != "" {
			fmt.Printf("\nSummary: %s\n", r.Summary)
		}
		if len(r.Highlights) > 0 {
			fmt.Println("\nHighlights:")
			for _, h := range r.Highlights {
				fmt.Printf("  • %s\n", h)
			}
		}
	}

	return nil
}
