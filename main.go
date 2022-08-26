package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type filePart struct {
	name       string
	start, end int
}

func downloader(url string) error {
	if strings.TrimSpace(url) == "" {
		return errors.New("invalid url")
	}

	res, err := http.Head(url)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return errors.New("unsupported protocol scheme")
	}
	urlSplit := strings.Split(url, "/")
	filename := urlSplit[len(urlSplit)-1]
	if res.Header.Get("Accept-Ranges") != "bytes" {
		return errors.New("impossible de télécharger ce fichier")
	}

	cntLen, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	nbPart := 3
	offset := cntLen / nbPart

	jobs := make(chan filePart, nbPart)
	results := make(chan string, nbPart)

	for w := 0; w < nbPart; w++ {
		go worker(url, jobs, results)
	}

	for i := 0; i < nbPart; i++ {
		name := fmt.Sprintf("part%d", i)
		start := i * offset
		end := (i + 1) * offset
		jobs <- filePart{name: name, start: start, end: end}
	}
	close(jobs)

	for i := 0; i < nbPart; i++ {
		fmt.Println(<-results)
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	for i := 0; i < nbPart; i++ {
		name := fmt.Sprintf("part%d", i)
		file, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		out.WriteAt(file, int64(i*offset))

		if err := os.Remove(name); err != nil {
			return err
		}
	}

	return nil
}

func worker(url string, jobs <-chan filePart, results chan<- string) {
	for job := range jobs {
		part, err := os.Create(job.name)
		if err != nil {
			log.Fatal(err)
		}
		client := http.Client{}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", job.start, job.end))
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.OpenFile(job.name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		bar := progressbar.DefaultBytes(
			res.ContentLength,
			"downloading",
		)
		io.Copy(io.MultiWriter(f, bar), res.Body)

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		_, err = part.Write(body)
		if err != nil {
			log.Fatal(err)
		}

		results <- part.Name()

		f.Close()
		res.Body.Close()
		part.Close()
	}
}

func main() {
	var url string
	flag.StringVar(&url, "u", "https://agritrop.cirad.fr/584726/1/Rapport.pdf", "url of the file to download")
	flag.Parse()
	start := time.Now()
	err := downloader(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(time.Since(start))
}
