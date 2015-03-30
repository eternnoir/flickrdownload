package flickrdownloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type FlickrDownloader struct {
	DebugMode   bool
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
	FatalLogger *log.Logger
}

func InitDownloader(debug bool) *FlickrDownloader {
	downloader := new(FlickrDownloader)
	downloader.DebugMode = debug
	downloader.InitLogger(os.Stdout, os.Stdout, os.Stderr)
	return downloader
}

func (downloader *FlickrDownloader) InitLogger(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	downloader.InfoLogger = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime)
	if downloader.DebugMode {
		downloader.DebugLogger = log.New(warningHandle,
			"DEBUG: ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}
	downloader.FatalLogger = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func (downloader *FlickrDownloader) SaveAllPhoto(url, path string, maxPage int) {
	pageUrls, err := downloader.getPagesUrls(url)
	if err != nil {
		downloader.InfoLogger.Fatal(err)
		return
	}
	photoPageUrls := []string{}
	for pageIndex, element := range pageUrls {
		if (pageIndex + 1) > maxPage {
			break
		}
		us, err := downloader.getPhotoUrls(element)
		if err != nil {
			downloader.InfoLogger.Fatal(err)
		}
		photoPageUrls = append(photoPageUrls, us...)
	}
	downloader.InfoLogger.Println("Finded " + strconv.Itoa(len(photoPageUrls)) + " photos.In " + url)
	var wg sync.WaitGroup
	for _, photoUrl := range photoPageUrls {
		wg.Add(1)
		go downloader.savePhoto(photoUrl, path, &wg)
		time.Sleep(time.Second * 1)
	}
	wg.Wait()
}

func (downloader *FlickrDownloader) getPagesUrls(url string) (uris []string, err error) {
	downloader.InfoLogger.Println("Find Page Urls " + url)
	return findAllPages(url)
}

func (downloader *FlickrDownloader) getPhotoUrls(url string) (uris []string, err error) {
	downloader.InfoLogger.Println("Find Photo Urls " + url)
	return findPhotoUrls(url)
}

func (downloader *FlickrDownloader) savePhoto(url, path string, wg *sync.WaitGroup) {
	trueLink, _ := findPhotoTrueLink(url, "o")
	downloader.InfoLogger.Println("Download " + trueLink)
	resp, err := http.Get(trueLink)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		downloader.InfoLogger.Fatal(err)
		return
	}
	filename := parseFileName(trueLink)
	downloader.InfoLogger.Println("Save " + filename)
	out, err := os.Create(path + "/" + filename)
	defer out.Close()
	if err != nil {
		downloader.InfoLogger.Fatal(err)
		return
	}
	_, ferr := io.Copy(out, resp.Body)
	if ferr != nil {
		downloader.InfoLogger.Fatal(ferr)
		return
	}
	downloader.InfoLogger.Println("File :" + filename + " Saved.")
	wg.Done()
}
