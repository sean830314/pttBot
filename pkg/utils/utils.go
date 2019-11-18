package utils

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	UserResponse = "請輸入:\n1. 文章版名\n\t 得此版最新一頁文章標題和文章ID\n2. 文章版名#文章ID\n\t 得此文章內容"
	NotFindBoard = "is not found board"
)

func GetDocument(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", "over18=1")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	return doc, err
}
