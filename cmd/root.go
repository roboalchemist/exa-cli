package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/roboalchemist/exa-cli/pkg/api"
	"github.com/roboalchemist/exa-cli/pkg/auth"
	"github.com/roboalchemist/exa-cli/pkg/output"
	"github.com/spf13/cobra"
)

var appVersion = "dev"

// Global flag values
var (
	flagJSON      bool
	flagPlaintext bool
	flagNoColor   bool
	flagDebug     bool
	flagFields    string
	flagJQ        string
)

var rootCmd = &cobra.Command{
	Use:   "exa",
	Short: "CLI for the Exa AI search API",
	Long: `exa is a command-line interface for the Exa AI search API.

Search the web, find similar pages, get AI-powered answers,
retrieve page contents, and explore code context.

Authentication:
  Set EXA_API_KEY environment variable or run 'exa auth'.

Examples:
  exa search "hottest AI startups"
  exa search "climate change research" --type deep -n 5
  exa answer "What is the capital of France?"
  exa similar "https://arxiv.org/abs/2307.06435"
  exa contents https://example.com
  exa context "React hooks state management"`,
	Version:       appVersion,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.BoolVarP(&flagJSON, "json", "j", false, "JSON output")
	pf.BoolVarP(&flagPlaintext, "plaintext", "p", false, "Tab-separated output for piping")
	pf.BoolVar(&flagNoColor, "no-color", false, "Disable colored output")
	pf.BoolVar(&flagDebug, "debug", false, "Verbose logging to stderr")
	pf.StringVar(&flagFields, "fields", "", "Comma-separated fields for JSON output")
	pf.StringVar(&flagJQ, "jq", "", "JQ expression to filter JSON output")
}

// GetOutputOptions builds output.Options from global flags.
func GetOutputOptions() output.Options {
	opts := output.Options{
		NoColor: flagNoColor,
		Debug:   flagDebug,
		Fields:  flagFields,
		JQ:      flagJQ,
	}
	switch {
	case flagJSON:
		opts.Mode = output.ModeJSON
	case flagPlaintext:
		opts.Mode = output.ModePlaintext
	default:
		opts.Mode = output.ModeTable
	}
	return opts
}

// newClient creates an authenticated API client.
func newClient() (*api.Client, error) {
	apiKey, err := auth.GetAPIKey()
	if err != nil {
		return nil, err
	}

	client := api.NewClient(auth.GetBaseURL(), apiKey)
	if flagDebug {
		client.SetDebug(DebugLog)
	}
	return client, nil
}

// newContext returns a background context.
func newContext() context.Context {
	return context.Background()
}

// Execute runs the root command.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		opts := GetOutputOptions()
		output.RenderError(err, opts)
	}
	return err
}

// SetVersion sets the application version.
func SetVersion(v string) {
	appVersion = v
	rootCmd.Version = v
	api.Version = v
}

// GetRootCmd returns the root command (used by cmd/gendocs).
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// DebugLog prints debug output to stderr.
func DebugLog(format string, args ...interface{}) {
	if flagDebug {
		fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
	}
}
