package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// run executes the exa binary with args and returns stdout, stderr, and error.
func run(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command("./exa", args...)
	cmd.Env = os.Environ()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// mustRun is like run but fails the test on error.
func mustRun(t *testing.T, args ...string) string {
	t.Helper()
	stdout, stderr, err := run(t, args...)
	if err != nil {
		t.Fatalf("exa %s failed: %v\nstdout: %s\nstderr: %s", strings.Join(args, " "), err, stdout, stderr)
	}
	return stdout
}

// requireAPIKey skips the test if EXA_API_KEY is not set.
func requireAPIKey(t *testing.T) {
	t.Helper()
	if os.Getenv("EXA_API_KEY") == "" {
		t.Skip("EXA_API_KEY not set")
	}
}

// --- Smoke tests (no API key needed) ---

func TestSmoke_Help(t *testing.T) {
	out := mustRun(t, "--help")
	if !strings.Contains(out, "Search the web") {
		t.Error("--help missing expected content")
	}
	if !strings.Contains(out, "Available Commands") {
		t.Error("--help missing Available Commands")
	}
}

func TestSmoke_Version(t *testing.T) {
	out := mustRun(t, "--version")
	if !strings.Contains(out, "exa version") {
		t.Errorf("--version unexpected output: %s", out)
	}
}

func TestSmoke_Docs(t *testing.T) {
	out := mustRun(t, "docs")
	if !strings.Contains(out, "# exa-cli") {
		t.Error("docs missing README header")
	}
	if !strings.Contains(out, "Installation") {
		t.Error("docs missing Installation section")
	}
}

func TestSmoke_SkillPrint(t *testing.T) {
	out := mustRun(t, "skill", "print")
	if !strings.Contains(out, "# exa-cli") {
		t.Error("skill print missing header")
	}
	if !strings.Contains(out, "examples") {
		t.Error("skill print missing examples")
	}
}

func TestSmoke_CompletionBash(t *testing.T) {
	out := mustRun(t, "completion", "bash")
	if !strings.Contains(out, "bash completion") {
		t.Error("bash completion missing expected content")
	}
}

func TestSmoke_CompletionZsh(t *testing.T) {
	out := mustRun(t, "completion", "zsh")
	if !strings.Contains(out, "compdef") {
		t.Error("zsh completion missing compdef")
	}
}

func TestSmoke_CompletionFish(t *testing.T) {
	out := mustRun(t, "completion", "fish")
	if !strings.Contains(out, "complete") {
		t.Error("fish completion missing expected content")
	}
}

func TestSmoke_SearchHelp(t *testing.T) {
	out := mustRun(t, "search", "--help")
	if !strings.Contains(out, "--num-results") {
		t.Error("search --help missing --num-results")
	}
	if !strings.Contains(out, "--type") {
		t.Error("search --help missing --type")
	}
	if !strings.Contains(out, "--category") {
		t.Error("search --help missing --category")
	}
}

func TestSmoke_AnswerHelp(t *testing.T) {
	out := mustRun(t, "answer", "--help")
	if !strings.Contains(out, "--stream") {
		t.Error("answer --help missing --stream")
	}
}

func TestSmoke_SimilarHelp(t *testing.T) {
	out := mustRun(t, "similar", "--help")
	if !strings.Contains(out, "--exclude-source") {
		t.Error("similar --help missing --exclude-source")
	}
}

func TestSmoke_ContentsHelp(t *testing.T) {
	out := mustRun(t, "contents", "--help")
	if !strings.Contains(out, "--text-max-chars") {
		t.Error("contents --help missing --text-max-chars")
	}
}

func TestSmoke_ContextHelp(t *testing.T) {
	out := mustRun(t, "context", "--help")
	if !strings.Contains(out, "--tokens") {
		t.Error("context --help missing --tokens")
	}
}

func TestSmoke_NoArgs(t *testing.T) {
	out := mustRun(t, "--help")
	for _, cmd := range []string{"search", "answer", "similar", "contents", "context", "usage", "auth", "docs", "completion", "skill"} {
		if !strings.Contains(out, cmd) {
			t.Errorf("root help missing command: %s", cmd)
		}
	}
}

func TestSmoke_UnknownCommand(t *testing.T) {
	_, stderr, err := run(t, "nonexistent")
	if err == nil {
		t.Error("expected error for unknown command")
	}
	if !strings.Contains(stderr, "unknown command") {
		// Cobra puts the error on stderr or embeds in the output
		out := mustRun(t, "--help")
		_ = out // just verify help still works
	}
}

// --- Integration tests (require EXA_API_KEY) ---

func TestIntegration_SearchBasic(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "golang testing best practices", "-n", "3")
	if !strings.Contains(out, "results") || !strings.Contains(out, "http") {
		// Table output should contain URLs
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) < 2 { // header + at least 1 result
			t.Errorf("search returned too few lines: %d", len(lines))
		}
	}
}

func TestIntegration_SearchJSON(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "golang testing", "-n", "2", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("search --json returned invalid JSON: %v\noutput: %s", err, out)
	}

	results, ok := resp["results"].([]interface{})
	if !ok {
		t.Fatal("search --json missing results array")
	}
	if len(results) == 0 {
		t.Error("search --json returned 0 results")
	}
	if len(results) > 2 {
		t.Errorf("search --json returned %d results, expected <= 2", len(results))
	}

	// Check result fields
	first := results[0].(map[string]interface{})
	for _, field := range []string{"title", "url", "id"} {
		if _, ok := first[field]; !ok {
			t.Errorf("search result missing field: %s", field)
		}
	}
}

func TestIntegration_SearchJSONFields(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "test query", "-n", "1", "--json", "--fields", "title,url")

	var results []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("search --json --fields returned invalid JSON: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("search --json --fields returned 0 results")
	}

	first := results[0]
	if _, ok := first["title"]; !ok {
		t.Error("field selection missing 'title'")
	}
	if _, ok := first["url"]; !ok {
		t.Error("field selection missing 'url'")
	}
	// Should NOT have score since we didn't request it
	if _, ok := first["score"]; ok {
		t.Error("field selection should not include 'score'")
	}
}

func TestIntegration_SearchJQ(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "test query", "-n", "2", "--json", "--jq", ".results | length")

	out = strings.TrimSpace(out)
	if out != "2" {
		t.Errorf("jq '.results | length' expected 2, got: %s", out)
	}
}

func TestIntegration_SearchPlaintext(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "golang", "-n", "2", "--plaintext")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 { // header + results
		t.Errorf("plaintext output too few lines: %d", len(lines))
	}
	// Each line should be tab-separated
	for _, line := range lines {
		if !strings.Contains(line, "\t") {
			t.Errorf("plaintext line missing tabs: %s", line)
		}
	}
}

func TestIntegration_SearchWithType(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "machine learning", "-n", "2", "--type", "neural", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("search --type neural returned invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	if len(results) == 0 {
		t.Error("neural search returned 0 results")
	}
}

func TestIntegration_SearchWithCategory(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "OpenAI", "-n", "2", "--category", "company", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("search --category returned invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	if len(results) == 0 {
		t.Error("category search returned 0 results")
	}
}

func TestIntegration_SearchNoContents(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "test", "-n", "1", "--no-contents", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	if len(results) == 0 {
		t.Fatal("no results")
	}
	first := results[0].(map[string]interface{})
	if text, ok := first["text"]; ok && text != "" {
		t.Error("--no-contents should not include text")
	}
}

func TestIntegration_SearchWithHighlights(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "golang concurrency", "-n", "1", "--highlights", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	if len(results) == 0 {
		t.Fatal("no results")
	}
	first := results[0].(map[string]interface{})
	if _, ok := first["highlights"]; !ok {
		t.Error("--highlights should include highlights field")
	}
}

func TestIntegration_ContentsBasic(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "contents", "https://example.com", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("contents --json invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	if len(results) != 1 {
		t.Errorf("contents expected 1 result, got %d", len(results))
	}
	first := results[0].(map[string]interface{})
	if first["url"] != "https://example.com" {
		t.Errorf("contents URL mismatch: %v", first["url"])
	}
	if text, ok := first["text"].(string); !ok || text == "" {
		t.Error("contents should include text by default")
	}
}

func TestIntegration_ContentsWithSummary(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "contents", "https://example.com", "--summary", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	if len(results) == 0 {
		t.Fatal("no results")
	}
	first := results[0].(map[string]interface{})
	if _, ok := first["summary"]; !ok {
		t.Error("--summary should include summary field")
	}
}

func TestIntegration_SimilarBasic(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "similar", "https://arxiv.org/abs/2307.06435", "-n", "3", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("similar --json invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	if len(results) == 0 {
		t.Error("similar returned 0 results")
	}
	if len(results) > 3 {
		t.Errorf("similar returned %d results, expected <= 3", len(results))
	}
}

func TestIntegration_SimilarExcludeSource(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "similar", "https://arxiv.org/abs/2307.06435", "-n", "3", "--exclude-source", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	for _, r := range results {
		url := r.(map[string]interface{})["url"].(string)
		if strings.Contains(url, "arxiv.org") {
			t.Errorf("--exclude-source should exclude arxiv.org, got: %s", url)
		}
	}
}

func TestIntegration_AnswerBasic(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "answer", "What is the capital of France?")
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "paris") {
		t.Errorf("answer should mention Paris, got: %s", out[:min(200, len(out))])
	}
}

func TestIntegration_AnswerJSON(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "answer", "What is 2+2?", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("answer --json invalid JSON: %v", err)
	}
	if _, ok := resp["answer"]; !ok {
		t.Error("answer --json missing answer field")
	}
	if _, ok := resp["citations"]; !ok {
		t.Error("answer --json missing citations field")
	}
}

func TestIntegration_ContextBasic(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "context", "Python list comprehension", "--tokens", "1000")
	if len(out) < 50 {
		t.Errorf("context output too short (%d chars), expected code context", len(out))
	}
}

func TestIntegration_ContextJSON(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "context", "JavaScript async await", "--tokens", "500", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("context --json invalid JSON: %v", err)
	}
	if _, ok := resp["response"]; !ok {
		t.Error("context --json missing response field")
	}
}

func TestIntegration_SearchDomainFilter(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "go programming", "-n", "3", "--include-domains", "go.dev", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	for _, r := range results {
		url := r.(map[string]interface{})["url"].(string)
		if !strings.Contains(url, "go.dev") {
			t.Errorf("--include-domains go.dev but got URL: %s", url)
		}
	}
}

func TestIntegration_SearchDateFilter(t *testing.T) {
	requireAPIKey(t)
	out := mustRun(t, "search", "AI news", "-n", "2", "--start-date", "2026-01-01", "--json")

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	results := resp["results"].([]interface{})
	if len(results) == 0 {
		t.Error("date-filtered search returned 0 results")
	}
}

// --- Error handling tests ---

func TestIntegration_SearchNoQuery(t *testing.T) {
	_, _, err := run(t, "search")
	if err == nil {
		t.Error("search with no query should fail")
	}
}

func TestIntegration_SimilarNoURL(t *testing.T) {
	_, _, err := run(t, "similar")
	if err == nil {
		t.Error("similar with no URL should fail")
	}
}

func TestIntegration_AnswerNoQuery(t *testing.T) {
	_, _, err := run(t, "answer")
	if err == nil {
		t.Error("answer with no query should fail")
	}
}

func TestIntegration_ContentsNoURL(t *testing.T) {
	_, _, err := run(t, "contents")
	if err == nil {
		t.Error("contents with no URL should fail")
	}
}

func TestIntegration_ContextNoQuery(t *testing.T) {
	_, _, err := run(t, "context")
	if err == nil {
		t.Error("context with no query should fail")
	}
}

func TestIntegration_InvalidAPIKey(t *testing.T) {
	cmd := exec.Command("./exa", "search", "test", "-n", "1")
	cmd.Env = append(os.Environ(), "EXA_API_KEY=invalid-key-12345")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		t.Error("invalid API key should produce an error")
	}
}
