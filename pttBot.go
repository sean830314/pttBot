// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pttBot/pkg/api/ptt"
	"pttBot/pkg/utils"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					message.Text = strings.Trim(message.Text, " \t\n\r")
					if len(strings.Split(message.Text, "#")) == 2 {
						//https: //www.ptt.cc/bbs/Gossiping/M.1571466186.A.F50.html
						articleUrl := fmt.Sprintf("%v/bbs/%v/%v", ptt.PttPrefix, strings.Split(message.Text, "#")[0], strings.Split(message.Text, "#")[1])
						//articleUrl := ptt.PttPrefix + "/bbs/" + strings.Split(message.Text, "#")[0] + "/" + strings.Split(message.Text, "#")[1]
						err, article := ptt.GetArticle(articleUrl)
						response := ""
						if err != nil {
							log.Print(err)
							response = utils.UserResponse
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
								log.Print(err)
							}
						} else if article.Title == "" {
							response = utils.UserResponse
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
								log.Print(err)
							}
						} else {
							response = fmt.Sprintf("Title: %v \n", article.Title)
							response += fmt.Sprintf("Author: %v \n", article.Author)
							response += fmt.Sprintf("Date: %v \n", article.Date)
							response += fmt.Sprintf("Ip: %v \n", article.Ip)
							response += fmt.Sprintf("Content: %v \n", article.Content)
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
								log.Print(err)
							}
						}
					} else {
						err, boardArticles := ptt.GetTopArticle(message.Text)
						response := ""
						if err != nil {
							response = err.Error()
							response += "\n" + utils.UserResponse
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
								log.Print(err)
							}
							log.Print(err)
						}
						if ptt.CheckBoardArticlesIsNil(boardArticles) {
							response += "\n" + utils.UserResponse
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
								log.Print(err)
							}
						}
						for i, ba := range boardArticles {
							response += fmt.Sprintf("第 %v 篇:\n%v\n", strconv.Itoa(i+1), ba.Title)
							response += fmt.Sprintf("文章版名:  %v\n文章ID: %v\n", ba.BoardName, ba.Id)
						}
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
							log.Print(err)
						}
					}

				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
