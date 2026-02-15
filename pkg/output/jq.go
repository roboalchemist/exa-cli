package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itchyny/gojq"
)

// RunJQ applies a jq expression to data and prints results to stdout.
func RunJQ(data interface{}, expr string) error {
	query, err := gojq.Parse(expr)
	if err != nil {
		return fmt.Errorf("invalid jq expression: %w", err)
	}

	// Convert to generic interface via JSON round-trip for gojq compatibility
	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal for jq: %w", err)
	}

	var input interface{}
	if err := json.Unmarshal(raw, &input); err != nil {
		return fmt.Errorf("unmarshal for jq: %w", err)
	}

	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return fmt.Errorf("jq error: %w", err)
		}

		out, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal jq result: %w", err)
		}
		fmt.Fprintln(os.Stdout, string(out))
	}

	return nil
}
