package flickrdownloader

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strconv"
	"strings"
)

var FLICKR_SITE = "https://www.flickr.com/"

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
	fmt.Println(maxpage)
	for i := 1; i <= maxpage; i++ {
		oneurl := url + "/page" + strconv.Itoa(i)
		returls = append(returls, oneurl)
	}
	return returls, nil
}

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
	fmt.Println("Finded " + strconv.Itoa(len(urls)) + " photos.In " + url)
	return urls, nil

}

func findPhotoTrueLink(url, size string) (uri string, err error) {
	doc, err := goquery.NewDocument(url + "/sizes/" + size)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	trueurl := ""
	doc.Find("#allsizes-photo img").Each(func(i int, s *goquery.Selection) {
		tr, isFind := s.Attr("src")
		if isFind {
			fmt.Println("TRUELINK: " + trueurl)
			trueurl = tr
		} else {
			fmt.Println("Not Found")
		}
	})
	return trueurl, nil
}

func ParseFileName(urls string) string {
	fileURL, err := url.Parse(urls)

	if err != nil {
		panic(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	return fileName
}

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
