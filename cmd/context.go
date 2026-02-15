package cmd

import (
	"fmt"
	"strings"

	"github.com/roboalchemist/exa-cli/pkg/api"
	"github.com/roboalchemist/exa-cli/pkg/output"
	"github.com/spf13/cobra"
)

var contextTokens int

var contextCmd = &cobra.Command{
	Use:   "context [query]",
	Short: "Get code context from Exa Code",
	Long: `Search for code-related context using Exa Code.

Returns relevant code snippets, documentation, and context
for a programming query.

Examples:
  exa context "React hooks state management"
  exa context "Python async await patterns" --tokens 5000
  exa context "Go error handling best practices" --json`,
	Args: cobra.MinimumNArgs(1),
	RunE: runContext,
}

func init() {
	contextCmd.Flags().IntVar(&contextTokens, "tokens", 0, "Token limit for response (0=dynamic)")

	rootCmd.AddCommand(contextCmd)
}

func runContext(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	req := &api.ContextRequest{
		Query: strings.Join(args, " "),
	}
	if contextTokens > 0 {
		req.TokensNum = contextTokens
	} else {
		req.TokensNum = "dynamic"
	}

	resp, err := client.GetContext(newContext(), req)
	if err != nil {
		return err
	}

	opts := GetOutputOptions()

	if opts.Mode == output.ModeJSON {
		return output.RenderJSON(resp, opts)
	}

	// Print context directly
	fmt.Println(resp.Context)

	if cost := resp.GetCost(); cost != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "\nCost: $%.4f\n", cost.Total)
	}

	return nil
}
