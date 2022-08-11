package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"))

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)
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
					switch message.Text {
					case "承認する":
						richMemu := linebot.RichMenu{
							Size: linebot.RichMenuSize{
								Width:  2500,
								Height: 1686,
							},
							Selected:    false,
							Name:        "Nice richmenu",
							ChatBarText: "Tap here",
							Areas: []linebot.AreaDetail{
								{
									Bounds: linebot.RichMenuBounds{
										X:      0,
										Y:      0,
										Width:  2500,
										Height: 1686,
									},
									Action: linebot.RichMenuAction{
										Type: "postback",
										Data: "action=test&itemid=123",
									},
								},
							},
						}
						res, err := bot.CreateRichMenu(richMemu).Do()
						if err != nil {
							log.Print(err)
						}
						fmt.Println(res.RichMenuID)
						if _, err := bot.UploadRichMenuImage(res.RichMenuID, "./image/richmenu.png").Do(); err != nil {
							log.Print(err)
						}

						if event.Source.UserID != os.Getenv("USER_ID") {
							bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("承認に失敗しました")).Do()
						} else {
							if _, err := bot.LinkUserRichMenu(event.Source.UserID, res.RichMenuID).Do(); err != nil {
								log.Print(err)
							}
							bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("承認されました")).Do()
						}
					case "社内ツールを使う":
						replyMessage := fmt.Sprintf("user id is %s", event.Source.UserID)
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							log.Print(err)
						}
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf("sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
