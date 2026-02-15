package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/roboalchemist/exa-cli/pkg/auth"
	"github.com/roboalchemist/exa-cli/pkg/output"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Configure API key authentication",
	Long: `Store your Exa API key in a local config file.

The API key is stored in ~/.exa-auth.json with mode 0600.
You can also set the EXA_API_KEY environment variable instead.

Get your API key at: https://dashboard.exa.ai/api-keys

Examples:
  exa auth                       # Interactive setup
  EXA_API_KEY=xxx exa search ... # Use env var instead`,
	RunE: runAuth,
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func runAuth(cmd *cobra.Command, args []string) error {
	opts := GetOutputOptions()

	// Check if already authenticated
	if key, err := auth.GetAPIKey(); err == nil && key != "" {
		configPath, _ := auth.ConfigPath()
		if os.Getenv("EXA_API_KEY") != "" {
			fmt.Println("Currently authenticated via EXA_API_KEY environment variable.")
		} else {
			fmt.Printf("Currently authenticated via config file: %s\n", configPath)
		}
		fmt.Println()
	}

	// Check if stdin is interactive
	fi, _ := os.Stdin.Stat()
	if fi.Mode()&os.ModeCharDevice == 0 {
		return fmt.Errorf("--api-key flag or EXA_API_KEY env var required in non-interactive mode")
	}

	fmt.Print("Enter your Exa API key: ")
	reader := bufio.NewReader(os.Stdin)
	key, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	key = strings.TrimSpace(key)

	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	if err := auth.SaveAuth(auth.AuthConfig{APIKey: key}); err != nil {
		return fmt.Errorf("save auth: %w", err)
	}

	configPath, _ := auth.ConfigPath()
	output.Success(fmt.Sprintf("API key saved to %s", configPath), opts)
	return nil
}
