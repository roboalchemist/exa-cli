package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for exa.

To load completions:

Bash:
  $ source <(exa completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ exa completion bash > /etc/bash_completion.d/exa
  # macOS:
  $ exa completion bash > $(brew --prefix)/etc/bash_completion.d/exa

Zsh:
  $ source <(exa completion zsh)

  # To load completions for each session, execute once:
  $ exa completion zsh > "${fpath[1]}/_exa"

Fish:
  $ exa completion fish | source

  # To load completions for each session, execute once:
  $ exa completion fish > ~/.config/fish/completions/exa.fish

PowerShell:
  PS> exa completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletionV2(os.Stdout, true)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
