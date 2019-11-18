package ptt

import (
	"errors"
	"fmt"
	"pttBot/pkg/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type BoardArticle struct {
	Id        string
	BoardName string
	Title     string
	Url       string
}

func GetTopArticle(boardName string) (error, []*BoardArticle) {
	url := "https://www.ptt.cc/bbs/" + boardName + "/index.html"
	doc, err := utils.GetDocument(url)
	if err != nil {
		return err, nil
	}
	not_found := doc.Find("div.bbs-content").Text()
	if not_found == "404 - Not Found." {
		return errors.New(fmt.Sprintf("Error: url %s not found", url)), nil
	}
	main_content := doc.Find("div#main-container")
	boardArticles := make([]*BoardArticle, 0)
	main_content.Find("div.r-ent").Each(func(i int, s *goquery.Selection) {
		ba := BoardArticle{}
		title := strings.Trim(s.Find("div.title").Text(), ": \t\n\r")
		url, exists := s.Find("div.title").Find("a").Attr("href")
		if !exists {
			fmt.Println("url href is not exists")
		}
		ba.Title = title
		ba.Url = url
		urlStr := strings.Split(url, "/")
		ba.Id = urlStr[3]
		ba.BoardName = urlStr[2]
		boardArticles = append(boardArticles, &ba)
	})
	return nil, boardArticles
}

func CheckBoardArticlesIsNil(boardArticles []*BoardArticle) bool {
	for _, ba := range boardArticles {
		if ba == nil {
			return true
		}
	}
	return false
}
