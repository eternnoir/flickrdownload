package flickrdownloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type FlickrDownloader struct {
	targetUrl string
}

func InitDownloader(url string) *FlickrDownloader {
	downloader := new(FlickrDownloader)
	downloader.targetUrl = url
	return downloader
}

func (downloader *FlickrDownloader) SaveAllPhoto(url, path string, maxPage int) {
	pageUrls, err := downloader.getPagesUrls(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	photoPageUrls := []string{}
	for pageIndex, element := range pageUrls {
		if (pageIndex + 1) > maxPage {
			break
		}
		us, err := downloader.getPhotoUrls(element)
		if err != nil {
			fmt.Println(err)
		}
		photoPageUrls = append(photoPageUrls, us...)
	}
	fmt.Println("Finded " + strconv.Itoa(len(photoPageUrls)) + " photos.In " + url)
	var wg sync.WaitGroup
	for _, photoUrl := range photoPageUrls {
		wg.Add(1)
		go downloader.savePhoto(photoUrl, path)
		time.Sleep(time.Second * 1)
	}
	wg.Wait()
}

func (downloader *FlickrDownloader) getPagesUrls(url string) (uris []string, err error) {
	fmt.Println("Find Page Urls " + url)
	return findAllPages(url)
}

func (downloader *FlickrDownloader) getPhotoUrls(url string) (uris []string, err error) {
	fmt.Println("Find Photo Urls " + url)
	return findPhotoUrls(url)
}

func (downloader *FlickrDownloader) savePhoto(url, path string) {
	trueLink, _ := findPhotoTrueLink(url, "o")
	fmt.Println("Download " + trueLink)
	resp, err := http.Get(trueLink)
	if err != nil {
		fmt.Println(err)
		return
	}
	filename := parseFileName(trueLink)
	fmt.Println("Save " + filename)
	out, err := os.Create(path + "/" + filename)
	defer out.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, ferr := io.Copy(out, resp.Body)
	if ferr != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("File :" + filename + " Saved.")
}
