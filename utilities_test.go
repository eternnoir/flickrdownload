package flickrdownloader

import (
	"testing"
)

func TestParseFileName(t *testing.T) {
	uri := "https://farm9.staticflickr.com/8597/15825501054_88440b4577_o.png"
	fileName := parseFileName(uri)
	if fileName != "15825501054_88440b4577_o.png" {
		t.Error("ParseFileName Error ")
	} else {
		t.Log("TestParseFileName Success")
	}
}

func TestFindPhotoTrueLink(t *testing.T) {
	photoUrl := "https://www.flickr.com/photos/marksein/15825501054/"
	trueUrl, err := findPhotoTrueLink(photoUrl, "o")
	t.Log("TrueUrl:" + trueUrl)
	if err != nil {
		t.Error(err)
	}
	if trueUrl != "https://farm9.staticflickr.com/8597/15825501054_88440b4577_o.png" {
		t.Error("TrueLink not match")
	} else {
		t.Log("TestFindPhotoTrueLink test success")
	}

}

func TestFindPhotoId(t *testing.T) {
	photoUrl := "https://www.flickr.com/photos/marksein/15825501054/"
	photoId := parsePhotoId(photoUrl)
	if photoId != "marksein/15825501054" {
		t.Error(photoId)
	}
}
