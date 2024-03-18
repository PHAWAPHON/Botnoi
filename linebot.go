package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
)

type Cartoon struct {
	Title            string   `json:"title"`
	Year             int      `json:"year"`
	Creator          []string `json:"creator"`
	Rating           string   `json:"rating"`
	Genre            []string `json:"genre"`
	RuntimeInMinutes int      `json:"runtime_in_minutes"`
	Episodes         int      `json:"episodes"`
	Image            string   `json:"image"`
	ID               int      `json:"id"`
}

func getCartoonsData() ([]Cartoon, error) {
	var cartoons []Cartoon
	cartoonsURL := "https://api.sampleapis.com/cartoons/cartoons2D"

	resp, err := http.Get(cartoonsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&cartoons)
	if err != nil {
		return nil, err
	}

	return cartoons, nil
}

func main() {
	bot, err := linebot.New(
		"91e16db81caf7b978f7aeb1852beabd8",
		"tXHRYYJuC+bAf78Ku5eR0MTdxtMXQKtAMSnParOxZU+UMKJVxElp46NtvYRwtQn4v+uMRV7QySHn02gFSP0LZKre1S/Lltc3j/mZulet+8Bio7C8Dg+5HpolBZ8uhc8r9vXLQITb30i6mQZgrswodwdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.POST("/webhook", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.JSON(http.StatusBadRequest, nil)
			} else {
				c.JSON(http.StatusInternalServerError, nil)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if strings.HasPrefix(message.Text, "Details of ") {
						title := strings.TrimPrefix(message.Text, "Details of ")
						cartoons, err := getCartoonsData()
						if err != nil {
							log.Print(err)
							return
						}
						for _, cartoon := range cartoons {
							if cartoon.Title == title {
								details := fmt.Sprintf("Title: %s\nYear: %d\nRating: %s\nGenre: %v\nEpisodes: %d\nRuntime: %d minutes\n",
									cartoon.Title, cartoon.Year, cartoon.Rating, cartoon.Genre, cartoon.Episodes, cartoon.RuntimeInMinutes)
								if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(details), linebot.NewImageMessage(cartoon.Image, cartoon.Image)).Do(); err != nil {
									log.Print(err)
								}
								break
							}
						}
					} else {
						switch message.Text {
						case "cartoons":
							cartoons, err := getCartoonsData()
							if err != nil {
								log.Print(err)
								return
							}

							var columns []*linebot.CarouselColumn
							for i, cartoon := range cartoons {
								log.Println(cartoon.Image)
								if i == 0 || i == 5 {
									continue // ที่skip เพราะว่า image error ครับ
								}
								if i >= 12 {

									break
								}

								column := linebot.NewCarouselColumn(
									cartoon.Image, cartoon.Title, fmt.Sprintf("Genre: %v\nYear: %v\nRating: %v", cartoon.Genre, cartoon.Year, cartoon.Rating),
									linebot.NewMessageAction("Details", "Details of "+cartoon.Title),
								)
								columns = append(columns, column)
							}

							carousel := linebot.NewCarouselTemplate(columns...)
							template := linebot.NewTemplateMessage("Cartoons", carousel)
							if _, err := bot.ReplyMessage(event.ReplyToken, template).Do(); err != nil {
								log.Print(err)
							}

						case "help":
							quickReplyItems := []*linebot.QuickReplyButton{
								linebot.NewQuickReplyButton("", linebot.NewMessageAction("Cartoons", "cartoons")),
								linebot.NewQuickReplyButton("", linebot.NewMessageAction("Help", "help")),
							}
							quickReply := linebot.NewQuickReplyItems(quickReplyItems...)
							if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("You can use the following commands:\n- cartoons: Show cartoons list\n- help: Show this help message").WithQuickReplies(quickReply)).Do(); err != nil {
								log.Print(err)
							}
						case "สวัสดี":
							if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("สวัสดีครับ ยินดีต้อนรับสู่บอทน่าน")).Do(); err != nil {
								log.Print(err)
							}

						default:
							buttons := linebot.NewButtonsTemplate(
								"", "Welcome to Narn Cartoon Bot", "Choose an option",
								linebot.NewMessageAction("Show Cartoons", "cartoons"),
								linebot.NewMessageAction("Help", "help"),
							)
							if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Welcome", buttons)).Do(); err != nil {
								log.Print(err)
							}
						}
					}
				}
			}
		}
	})
	router.Run(":5678")
}
