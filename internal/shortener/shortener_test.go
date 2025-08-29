package shortener

import "testing"

func TestGenerateCode(t *testing.T) {
	code, err := GenerateCode(6)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 6 {
		t.Fatalf("want len 6, got %d", len(code))
	}
}

func TestNormalizeURL(t *testing.T) {
	u, err := NormalizeURL("https://example.com/path")
	if err != nil {
		t.Fatal(err)
	}
	if u[:5] != "https" {
		t.Fatalf("scheme should be https by default, got %s", u)
	}
}
