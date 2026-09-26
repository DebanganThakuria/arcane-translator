package service

import (
	"errors"
	"testing"
)

func TestPageHTMLUsesPastedHTML(t *testing.T) {
	pasted := "<html>pasted</html>"

	got, err := pageHTML("ixdzs", "https://ixdzs.tw/read/526058/p1.html", &pasted)
	if err != nil {
		t.Fatalf("pageHTML returned error: %v", err)
	}
	if got != pasted {
		t.Fatalf("pageHTML = %q, want the pasted HTML", got)
	}
}

func TestPageHTMLRefusesToScrapeManualSources(t *testing.T) {
	// An unroutable URL: if the refusal regressed, the test fails on the
	// error type rather than quietly reaching the real site.
	_, err := pageHTML("ixdzs", "http://127.0.0.1:1/read/526058/", nil)
	if !errors.Is(err, ErrManualHTMLRequired) {
		t.Fatalf("pageHTML error = %v, want ErrManualHTMLRequired", err)
	}
}
