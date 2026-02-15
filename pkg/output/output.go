package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Mode represents the output format.
type Mode int

const (
	ModeTable     Mode = iota // Default: colored table
	ModePlaintext             // Tab-separated, no colors
	ModeJSON                  // Pretty-printed JSON
)

// Options controls output rendering.
type Options struct {
	Mode    Mode
	NoColor bool
	Debug   bool
	Fields  string
	JQ      string
}

// TableData holds rows and headers for table rendering.
type TableData struct {
	Headers []string
	Rows    [][]string
	Footer  string
}

// RenderTable renders data in the appropriate output mode.
func RenderTable(td TableData, data interface{}, opts Options) error {
	switch opts.Mode {
	case ModeJSON:
		return renderJSONOutput(data, opts)
	case ModePlaintext:
		return renderPlaintext(os.Stdout, td)
	default:
		return renderTable(os.Stdout, td, opts)
	}
}

// RenderJSON outputs raw data as JSON.
func RenderJSON(data interface{}, opts Options) error {
	return renderJSONOutput(data, opts)
}

func renderJSONOutput(data interface{}, opts Options) error {
	data = FilterFields(data, opts.Fields)

	if opts.JQ != "" {
		return RunJQ(data, opts.JQ)
	}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	_, err = fmt.Fprintln(os.Stdout, string(out))
	return err
}

func renderPlaintext(w *os.File, td TableData) error {
	if len(td.Headers) > 0 {
		fmt.Fprintln(w, strings.Join(td.Headers, "\t"))
	}
	for _, row := range td.Rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	return nil
}

func renderTable(w *os.File, td TableData, opts Options) error {
	table := tablewriter.NewWriter(w)

	if shouldColor(opts) {
		colored := make([]string, len(td.Headers))
		for i, h := range td.Headers {
			colored[i] = color.New(color.FgCyan, color.Bold).Sprint(h)
		}
		table.SetHeader(colored)
	} else {
		table.SetHeader(td.Headers)
	}

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("   ")
	table.SetNoWhiteSpace(true)

	for _, row := range td.Rows {
		table.Append(row)
	}
	table.Render()

	if td.Footer != "" {
		if shouldColor(opts) {
			fmt.Fprintf(w, "\n%s\n", color.New(color.FgHiBlack).Sprint(td.Footer))
		} else {
			fmt.Fprintf(w, "\n%s\n", td.Footer)
		}
	}

	return nil
}

func shouldColor(opts Options) bool {
	if opts.NoColor {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// Error outputs an error message respecting the output mode.
func Error(message string, opts Options) {
	switch opts.Mode {
	case ModeJSON:
		_ = json.NewEncoder(os.Stderr).Encode(map[string]string{"error": message})
	default:
		fmt.Fprintf(os.Stderr, "%s %s\n", color.New(color.FgRed).Sprint("Error:"), message)
	}
}

// Success outputs a success message respecting the output mode.
func Success(message string, opts Options) {
	switch opts.Mode {
	case ModeJSON:
		_ = json.NewEncoder(os.Stderr).Encode(map[string]string{"status": "success", "message": message})
	default:
		fmt.Printf("%s %s\n", color.New(color.FgGreen).Sprint("OK:"), message)
	}
}

// PrintText prints text output (for answer/context commands).
func PrintText(text string) {
	fmt.Print(text)
}
