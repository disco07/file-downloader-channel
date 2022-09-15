package main

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestDownloader(t *testing.T) {
	var downloaderTests = []struct {
		description string
		url         string
		expected    error
	}{
		{
			description: "Download work",
			url:         "https://agritrop.cirad.fr/584726/1/Rapport.pdf",
			expected:    nil,
		},
		{
			description: "invalid url",
			url:         "",
			expected:    fmt.Errorf("invalid url"),
		},
		{
			description: "unable to download file with multithreading",
			url:         "https://youtu.be/w0NQlEMjntI",
			expected:    errors.New("unable to download file with multithreads"),
		},
		{
			description: "unable to parse variable",
			url:         "https://github.com/disco07/file-downloader",
			expected:    errors.New("unable to parse variable"),
		},
	}

	for _, tt := range downloaderTests {
		t.Run(tt.description, func(t *testing.T) {
			err := downloader(tt.url)
			if tt.expected == nil && err != nil {
				t.Errorf("Unexpected error for input %v: %v (expected %v)", tt.url, err, tt.expected)
			}

			if tt.expected != nil && err.Error() != tt.expected.Error() {
				t.Errorf("Unexpected error for input %v: %v (expected %v)", tt.url, err, tt.expected)
			}
		})
	}
}

func TestAppend(t *testing.T) {
	f, err := os.Create("filename")
	if err != nil {
		t.Errorf("failed to create file: %v", err.Error())
	}
	defer f.Close()

	part, err := os.Create("part0")
	if err != nil {
		t.Errorf("failed to create file: %v", err.Error())
	}

	_, err = os.Create("part1")
	if err != nil {
		t.Errorf("failed to create file: %v", err.Error())
	}

	tests := []struct {
		description string
		part        int
		offset      int
		file        *os.File
		expected    error
	}{
		{
			description: "append data in file",
			part:        0,
			offset:      1000,
			file:        f,
			expected:    nil,
		},
		//{
		//	description: "unable to delete file",
		//	part:        1,
		//	offset:      1000,
		//	file:        f,
		//	expected:    errors.New("cannot delete the file because it is still open"),
		//},
		{
			description: "file not found",
			part:        2,
			offset:      1000,
			file:        f,
			expected:    errors.New("file not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			part.Close()
			err := Append(tt.file, tt.part, tt.offset)
			if tt.expected == nil && err != nil {
				t.Errorf("got %v want %v", err.Error(), tt.expected)
			}

			if tt.expected != nil && err.Error() != tt.expected.Error() {
				t.Errorf("got %v want %v", err.Error(), tt.expected)
			}
		})
	}
}

func BenchmarkDownloader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		downloader("https://agritrop.cirad.fr/584726/1/Rapport.pdf")
	}
}
