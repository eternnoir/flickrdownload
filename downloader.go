package flickrdownloader

import (
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

// Create FlickerDownloader. Use debug para to setup logger in debug mode.
func InitDownloader(debug bool) *FlickrDownloader {
	downloader := new(FlickrDownloader)
	downloader.DebugMode = debug
	downloader.InitLogger(os.Stdout, os.Stdout, os.Stderr)
	return downloader
}

// Init logger. This method will init INFO,DEBUG,ERROR three logger to
// FlickerDownloader.
func (downloader *FlickrDownloader) InitLogger(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	downloader.InfoLogger = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime)
	downloader.DebugLogger = log.New(warningHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	downloader.FatalLogger = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// Save all photo by url. Here's url can contain manay page,it often be a set or
// an user's all photo.
// path is where you want to storage downloaded photo.
// A url often has many page, sometimes you dont want download all page at one times,
// you can use maxPage para.
// imageSize:
// 		o means origin.
// 		l means large.
//		m means Medium.
func (downloader *FlickrDownloader) SaveAllPhoto(url, path string, maxPage int, imageSize string) {
	pageUrls, err := downloader.getPagesUrls(url)
	if err != nil {
		downloader.errors(err)
		return
	}
	photoPageUrls := []string{}
	for pageIndex, element := range pageUrls {
		if (pageIndex + 1) > maxPage {
			break
		}
		us, err := downloader.getPhotoUrls(element)
		if err != nil {
			downloader.errors(err)
		}
		photoPageUrls = append(photoPageUrls, us...)
	}
	downloader.info("Finded " + strconv.Itoa(len(photoPageUrls)) + " photos.In " + url)
	var wg sync.WaitGroup
	for _, photoUrl := range photoPageUrls {
		wg.Add(1)
		go downloader.savePhoto(photoUrl, path, imageSize, &wg)
		// wait for a little time.
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
}

func (downloader *FlickrDownloader) getPagesUrls(url string) (uris []string, err error) {
	downloader.info("Find Page Urls " + url)
	return findAllPages(url)
}

func (downloader *FlickrDownloader) getPhotoUrls(url string) (uris []string, err error) {
	downloader.info("Find Photo Urls " + url)
	return findPhotoUrls(url)
}

func (downloader *FlickrDownloader) savePhoto(url, path, imageSize string, wg *sync.WaitGroup) {
	trueLink, _ := findPhotoTrueLink(url, imageSize)
	downloader.debug("Download " + trueLink)
	resp, err := http.Get(trueLink)
	defer resp.Body.Close()
	if err != nil {
		downloader.errors(err)
		return
	}
	filename := parseFileName(trueLink)
	downloader.debug("Save " + filename)
	out, err := os.Create(path + "/" + filename)
	defer out.Close()
	if err != nil {
		downloader.errors(err)
		return
	}
	_, ferr := io.Copy(out, resp.Body)
	if ferr != nil {
		downloader.errors(ferr)
		return
	}
	downloader.info("File :" + filename + " Saved.")
	wg.Done()
}

func (downloader *FlickrDownloader) info(v ...interface{}) {
	go downloader.InfoLogger.Println(v)
}
func (downloader *FlickrDownloader) debug(v ...interface{}) {
	if downloader.DebugMode {
		go downloader.InfoLogger.Println(v)
	}
}
func (downloader *FlickrDownloader) errors(v ...interface{}) {
	go downloader.InfoLogger.Println(v)
}
