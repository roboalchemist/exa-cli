package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// CLIError is a structured error for JSON output.
type CLIError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Recoverable bool   `json:"recoverable"`
	Suggestion  string `json:"suggestion,omitempty"`
}

// RenderError outputs an error in the appropriate format.
func RenderError(err error, opts Options) {
	if opts.Mode == ModeJSON {
		cliErr := toStructuredError(err)
		_ = json.NewEncoder(os.Stderr).Encode(cliErr)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func toStructuredError(err error) CLIError {
	msg := err.Error()

	switch {
	case strings.Contains(msg, "EXA_API_KEY"):
		return CLIError{
			Code:        "AUTH_REQUIRED",
			Message:     msg,
			Recoverable: true,
			Suggestion:  "Set EXA_API_KEY environment variable or run 'exa auth'",
		}
	case strings.Contains(msg, "status 401"):
		return CLIError{
			Code:        "AUTH_INVALID",
			Message:     "Invalid API key",
			Recoverable: true,
			Suggestion:  "Check your EXA_API_KEY value",
		}
	case strings.Contains(msg, "status 429"):
		return CLIError{
			Code:        "RATE_LIMITED",
			Message:     msg,
			Recoverable: true,
			Suggestion:  "Wait and retry, or reduce request frequency",
		}
	case strings.Contains(msg, "request failed"):
		return CLIError{
			Code:        "NETWORK_ERROR",
			Message:     msg,
			Recoverable: true,
			Suggestion:  "Check network connectivity",
		}
	default:
		return CLIError{
			Code:    "UNKNOWN",
			Message: msg,
		}
	}
}
