package main

import "testing"

func BenchmarkDownloader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		downloader("https://agritrop.cirad.fr/584726/1/Rapport.pdf")
	}
}
