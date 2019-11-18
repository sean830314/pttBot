package ptt

import (
	"errors"
	"fmt"
	"pttBot/pkg/utils"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	PttPrefix = "https://www.ptt.cc"
)

type Comment struct {
	PushTag        string
	PushUserID     string
	PushContent    string
	PushIpdatetime string
}
type Article struct {
	Title               string
	Author              string
	Date                string
	Content             string
	Ip                  string
	Comments            []Comment
	All, Count, P, B, N int
}

func GetArticle(url string) (error, *Article) {
	doc, err := utils.GetDocument(url)
	if err != nil {
		return err, nil
	}
	not_found := doc.Find("div.bbs-content").Text()
	if not_found == "404 - Not Found." {
		return errors.New(fmt.Sprintf("Error: url %s not found", url)), nil
	}
	main_content := doc.Find("div#main-content")
	article := &Article{}
	article.getArticleMetaline(main_content)
	article.getArticlePushTag(main_content)
	err = article.getArticleIP(main_content)
	if err != nil {
		return err, nil
	}
	err = article.getArticleContent(main_content)
	if err != nil {
		return err, nil
	}
	return nil, article
}

func (a *Article) getArticleMetaline(content *goquery.Selection) {
	// get meta
	content.Find("div.article-metaline").Each(func(i int, s *goquery.Selection) {
		k := s.Find("span.article-meta-tag").Text()
		v := s.Find("span.article-meta-value").Text()
		switch k {
		case "作者":
			a.Author = v
		case "標題":
			a.Title = v
		case "時間":
			a.Date = v
		}
		// remove article metaline
		s.Remove()
	})
	// remove article metaline-right
	content.Find("div.article-metaline-right").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})
}

func (a *Article) getArticlePushTag(content *goquery.Selection) {
	// get comments
	pushes := content.Find("div.push")
	a.Comments = make([]Comment, pushes.Size())
	pushes.Each(func(i int, push *goquery.Selection) {
		push_tag := strings.Trim(push.Find("span.push-tag").Text(), " \t\n\r")
		push_user_id := strings.Trim(push.Find("span.push-userid").Text(), " \t\n\r")
		push_content := strings.Trim(push.Find("span.push-content").Text(), ": \t\n\r")
		push_ipdatetime := strings.Trim(push.Find("span.push-ipdatetime").Text(), " \t\n\r")
		switch push_tag {
		case "推":
			a.P += 1
		case "噓":
			a.B += 1
		default:
			a.N += 1
		}
		a.Comments[i] = Comment{PushTag: push_tag,
			PushUserID:     push_user_id,
			PushContent:    push_content,
			PushIpdatetime: push_ipdatetime}
		push.Remove()
	})
	a.All = a.P + a.B + a.N
	a.Count = a.P - a.B
}

func (a *Article) getArticleIP(content *goquery.Selection) error {
	// get ip
	html, err := content.Html()
	if err != nil {
		return err
	}
	r, err := regexp.Compile("(※ 發信站: ).*")
	if err != nil {
		return err
	}
	ip := r.FindString(html)
	r, err = regexp.Compile("[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+")
	if err != nil {
		return err
	}
	a.Ip = r.FindString(ip)
	return nil
}

func (a *Article) getArticleContent(content *goquery.Selection) error {
	// remove redundant f2 class and remain text of others class
	content.Find("*").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "※ 發信站:") || strings.Contains(text, "※ 文章網址:") || strings.Contains(text, "※ 編輯:") {
			s.Remove()
		} else {
			s.ReplaceWithHtml(text)
		}
	})
	articleContent, err := content.Html()
	if err != nil {
		return err
	}
	a.Content = strings.Trim(articleContent, "-\t\n\r")
	return nil
}
