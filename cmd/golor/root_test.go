package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestConvertHSL(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := execute(&stdout, &stderr, []string{"convert", "#ff6b35", "--to", "hsl"}); err != nil {
		t.Fatalf("execute() error = %v, stderr = %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "HSL{H:16.") {
		t.Fatalf("convert output = %q", out)
	}
}

func TestConvertJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := execute(&stdout, &stderr, []string{"--json", "convert", "#ff6b35", "--to", "cmyk"}); err != nil {
		t.Fatalf("execute() error = %v, stderr = %s", err, stderr.String())
	}
	var got map[string]float64
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; output = %s", err, stdout.String())
	}
	if _, ok := got["c"]; !ok {
		t.Fatalf("expected lowercase cmyk JSON keys, got %v", got)
	}
}

func TestContrastJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := execute(&stdout, &stderr, []string{"contrast", "#ffffff", "#000000", "--json"}); err != nil {
		t.Fatalf("execute() error = %v, stderr = %s", err, stderr.String())
	}
	var got struct {
		Ratio float64 `json:"ratio"`
		AA    bool    `json:"aa"`
		AAA   bool    `json:"aaa"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; output = %s", err, stdout.String())
	}
	if got.Ratio != 21 || !got.AA || !got.AAA {
		t.Fatalf("contrast JSON = %+v", got)
	}
}

func TestPreviewColorModes(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := execute(&stdout, &stderr, []string{"preview", "#ff0000", "--color-mode", "truecolor"}); err != nil {
		t.Fatalf("execute() error = %v, stderr = %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "\x1b[48;2;255;0;0m") {
		t.Fatalf("truecolor preview output = %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if err := execute(&stdout, &stderr, []string{"preview", "#ff0000", "--color-mode", "256"}); err != nil {
		t.Fatalf("execute() error = %v, stderr = %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "\x1b[48;5;9m") {
		t.Fatalf("256 preview output = %q", stdout.String())
	}
}

func TestPaletteMissingFileJSONError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := execute(&stdout, &stderr, []string{"--json", "palette", "missing.png", "-n", "5"})
	if err == nil {
		t.Fatal("expected missing file error")
	}
	var got map[string]string
	if jsonErr := json.Unmarshal(stderr.Bytes(), &got); jsonErr != nil {
		t.Fatalf("json.Unmarshal() error = %v; stderr = %s", jsonErr, stderr.String())
	}
	if got["error"] == "" {
		t.Fatalf("expected structured error, got %v", got)
	}
}
