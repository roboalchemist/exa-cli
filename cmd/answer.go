package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/roboalchemist/exa-cli/pkg/api"
	"github.com/roboalchemist/exa-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	answerStream       bool
	answerText         bool
	answerOutputSchema string
)

var answerCmd = &cobra.Command{
	Use:   "answer [query]",
	Short: "Get an AI-powered answer with citations",
	Long: `Get an LLM-generated answer to your question with cited sources.

The answer is generated using Exa's search results as context,
providing grounded, factual responses with source citations.

Examples:
  exa answer "What is the capital of France?"
  exa answer "What are the latest AI breakthroughs in 2025?"
  exa answer "How does photosynthesis work?" --text
  exa answer "Explain quantum computing" --stream
  exa answer "List top 5 programming languages" --json`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAnswer,
}

func init() {
	f := answerCmd.Flags()
	f.BoolVar(&answerStream, "stream", false, "Stream the answer")
	f.BoolVar(&answerText, "text", false, "Include full text in citations")
	f.StringVar(&answerOutputSchema, "output-schema", "", "JSON schema file for structured output")

	rootCmd.AddCommand(answerCmd)
}

func runAnswer(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	req := &api.AnswerRequest{
		Query: strings.Join(args, " "),
		Text:  answerText,
	}

	if answerOutputSchema != "" {
		data, err := os.ReadFile(answerOutputSchema)
		if err != nil {
			return fmt.Errorf("read schema file: %w", err)
		}
		var schema interface{}
		if err := json.Unmarshal(data, &schema); err != nil {
			return fmt.Errorf("parse schema: %w", err)
		}
		req.OutputSchema = schema
	}

	opts := GetOutputOptions()

	// Streaming mode
	if answerStream && opts.Mode != output.ModeJSON {
		var finalResp *api.AnswerResponse
		err := client.AnswerStream(newContext(), req,
			func(text string) {
				fmt.Print(text)
			},
			func(resp *api.AnswerResponse) {
				finalResp = resp
			},
		)
		if err != nil {
			return err
		}

		fmt.Println() // Newline after streamed answer

		// Print citations
		if finalResp != nil && len(finalResp.Citations) > 0 {
			fmt.Println()
			cyan := color.New(color.FgCyan)
			cyan.Println("Sources:")
			for i, c := range finalResp.Citations {
				fmt.Printf("  %d. %s — %s\n", i+1, c.Title, c.URL)
			}
			if finalResp.CostDollars != nil {
				fmt.Printf("\nCost: $%.4f\n", finalResp.CostDollars.Total)
			}
		}

		return nil
	}

	// Non-streaming
	resp, err := client.Answer(newContext(), req)
	if err != nil {
		return err
	}

	if opts.Mode == output.ModeJSON {
		return output.RenderJSON(resp, opts)
	}

	// Pretty print answer
	fmt.Println(resp.Answer)

	if len(resp.Citations) > 0 {
		fmt.Println()
		cyan := color.New(color.FgCyan)
		cyan.Println("Sources:")
		for i, c := range resp.Citations {
			fmt.Printf("  %d. %s — %s\n", i+1, c.Title, c.URL)
		}
	}

	if resp.CostDollars != nil {
		fmt.Printf("\nCost: $%.4f\n", resp.CostDollars.Total)
	}

	return nil
}
