package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"

func FetchHtml(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		res.Body.Close()
		return nil, errors.New("Expected Content-Type to be text/html")
	}

	return res, nil
}

type PageData struct {
	Data        []byte
	ContentType string
}

func FetchPage(chapterUrl, url string) (PageData, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PageData{}, err
	}

	req.Header.Set("Referer", chapterUrl)
	req.Header.Set("User-Agent", UserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return PageData{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return PageData{}, errors.New("Returned non success status code")
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PageData{}, err
	}

	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		return PageData{}, errors.New("No Content Type")
	}

	return PageData{
		Data:        data,
		ContentType: contentType,
	}, nil
}

func GetExtention(contentType string) (string, error) {
	switch contentType {
	case "image/jpeg":
		return ".jpeg", nil
	case "image/png":
		return ".png", nil
	}

	return "", fmt.Errorf("Unknown Content-Type: %v", contentType)
}
