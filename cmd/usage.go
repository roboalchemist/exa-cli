package cmd

import (
	"fmt"
	"time"

	"github.com/roboalchemist/exa-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	usageStartDate string
	usageEndDate   string
	usageKeyID     string
)

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show API usage and costs",
	Long: `Display API usage statistics and costs for your account.

Shows request counts and credit usage over a time period.

Examples:
  exa usage
  exa usage --start-date 2025-01-01 --end-date 2025-01-31
  exa usage --json`,
	RunE: runUsage,
}

func init() {
	f := usageCmd.Flags()
	f.StringVar(&usageStartDate, "start-date", "", "Start of period (default: 30 days ago)")
	f.StringVar(&usageEndDate, "end-date", "", "End of period (default: now)")
	f.StringVar(&usageKeyID, "key-id", "", "Specific API key ID")

	rootCmd.AddCommand(usageCmd)
}

func runUsage(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	// Default date range: last 30 days
	now := time.Now()
	if usageEndDate == "" {
		usageEndDate = now.Format("2006-01-02")
	}
	if usageStartDate == "" {
		usageStartDate = now.AddDate(0, 0, -30).Format("2006-01-02")
	}

	// If no key ID specified, list keys and use the first one
	keyID := usageKeyID
	if keyID == "" {
		keys, err := client.ListAPIKeys(newContext())
		if err != nil {
			return fmt.Errorf("list API keys: %w", err)
		}
		if len(keys.APIKeys) == 0 {
			return fmt.Errorf("no API keys found")
		}
		keyID = keys.APIKeys[0].ID
		DebugLog("Using API key: %s (%s)", keyID, keys.APIKeys[0].Name)
	}

	resp, err := client.GetUsage(newContext(), keyID, usageStartDate, usageEndDate)
	if err != nil {
		return err
	}

	opts := GetOutputOptions()

	if opts.Mode == output.ModeJSON {
		return output.RenderJSON(resp, opts)
	}

	td := output.TableData{
		Headers: []string{"DATE", "REQUESTS", "CREDITS"},
	}

	totalRequests := 0
	totalCredits := 0.0
	for _, u := range resp.Usage {
		td.Rows = append(td.Rows, []string{
			u.Date,
			fmt.Sprintf("%d", u.RequestCount),
			fmt.Sprintf("%.4f", u.CreditUsage),
		})
		totalRequests += u.RequestCount
		totalCredits += u.CreditUsage
	}

	td.Footer = fmt.Sprintf("Total: %d requests, %.4f credits | %s to %s",
		totalRequests, totalCredits, usageStartDate, usageEndDate)

	return output.RenderTable(td, resp, opts)
}
