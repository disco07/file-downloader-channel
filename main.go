package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func downloader(url string) error {
	if strings.TrimSpace(url) == "" {
		return errors.New("invalid url")
	}

	ch := make(chan os.File)
	//defer close(ch)

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
	nbPart := 5
	offset := cntLen / nbPart

	for i := 0; i < nbPart; i++ {
		name := fmt.Sprintf("part%d", i)
		start := i * offset
		end := (i + 1) * offset
		part, err := os.Create(name)
		if err != nil {
			return err
		}

		go func() {
			ch <- writeFile(url, start, end, *part)
		}()
		<-ch
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

func writeFile(url string, start, end int, file os.File) os.File {
	defer file.Close()
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(body)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func main() {
	var url string
	flag.StringVar(&url, "u", "https://agritrop.cirad.fr/584726/1/Rapport.pdf", "url of the file to download")
	flag.Parse()
	err := downloader(url)
	if err != nil {
		log.Fatal(err)
		return
	}
}
