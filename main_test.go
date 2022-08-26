package main

import (
	"errors"
	"fmt"
	"testing"
)

var downloaderTests = []struct {
	url      string
	expected error
}{
	//{"https://agritrop.cirad.fr/584726/1/Rapport.pdf", nil},
	{"", fmt.Errorf("invalid url")},
	{"https://youtu.be/w0NQlEMjntI", errors.New("unable to download file with multithreads")},
	{"https://github.com/disco07/file-downloader", errors.New("unable to parse variable")},
}

func TestDownloader(t *testing.T) {
	for _, tt := range downloaderTests {
		err := downloader(tt.url)
		if err == nil {
			t.Errorf("Expected error for input %v", tt.url)
		}
		if err.Error() != tt.expected.Error() {
			t.Errorf("Unexpected error for input %v: %v (expected %v)", tt.url, err, tt.expected)
		}
	}
}

func BenchmarkDownloader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		downloader("https://agritrop.cirad.fr/584726/1/Rapport.pdf")
	}
}
