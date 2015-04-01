package flickrdownloader

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strconv"
	"strings"
)

var FLICKR_SITE = "https://www.flickr.com"

// Find all pages url by user or set.
func findAllPages(url string) (urls []string, err error) {
	maxpage := 0
	returls := []string{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return nil, err
	}

	doc.Find(".Paginator .pages .rapidnofollow").Each(func(i int, s *goquery.Selection) {
		num, err := strconv.Atoi(s.Text())
		if err == nil {
			if num > maxpage {
				maxpage = num
			}
		}
	})
	if maxpage == 0 {
		returls = append(returls, url+"/page1")
	}
	for i := 1; i <= maxpage; i++ {
		oneurl := url + "/page" + strconv.Itoa(i)
		returls = append(returls, oneurl)
	}
	return returls, nil
}

// Find all photo url from page.
func findPhotoUrls(url string) (uris []string, err error) {
	urls := []string{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".photo-display-item .hover-target .thumb .photo_container .rapidnofollow").Each(func(i int, s *goquery.Selection) {
		t, _ := s.Attr("href")
		urls = append(urls, FLICKR_SITE+t)
	})
	return urls, nil

}

// Find photo .jpg link by photo url.
// It will depance on size.
func findPhotoTrueLink(url, size string) (uri string, err error) {
	photoId := parsePhotoId(url)
	fmt.Println(FLICKR_SITE + "/" + photoId + "/sizes/" + size)
	doc, err := goquery.NewDocument(FLICKR_SITE + "/" + photoId + "/sizes/" + size)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	trueurl := ""
	doc.Find("#allsizes-photo img").Each(func(i int, s *goquery.Selection) {
		tr, isFind := s.Attr("src")
		if isFind {
			trueurl = tr
		}
	})
	return trueurl, nil
}

// Parse photo id by photo link.
// ex; https://www.flickr.com/photos/marksein/9448406987/in/set-72157634949960809
// will return marksein/9448406987
// ex: https://www.flickr.com/photos/marksein/9448406987
// will return marksein/9448406987 too.
func parsePhotoId(urls string) string {
	fmt.Println(urls)
	fileURL, err := url.Parse(urls)
	if err != nil {
		panic(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	id := segments[2] + "/" + segments[3]
	return id
}

// Get Filename by url.
// ex: http://example.com/ex.jpg
// It will return ex.jpg
func parseFileName(urls string) string {
	fileURL, err := url.Parse(urls)

	if err != nil {
		panic(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	return fileName
}
