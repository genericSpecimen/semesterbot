package main

import (
	"log"
	"os"
	"fmt"
	//"time"
	"strings"
	"net/http"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func getNotices(msg *tgba.MessageConfig) {
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	// Request the HTML page.
	base_url := "https://rlacollege.edu.in/"
	res, err := http.Get(base_url + "view-all-details.php")
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		msg.Text = "Failed to get notices; is the website working?"
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#news a").Each(func(i int, s *goquery.Selection) {
		notice_url, _ := s.Attr("href")

		parts := strings.Split(notice_url, "/")
		// ) is wrongfully parsed as markdown syntax, resulting in broken links in msg
		parts[2] = url.QueryEscape(parts[2])
		notice_url = strings.Join(parts[:], "/")

		title := s.Text()
		msg.Text += fmt.Sprintf("%d: [%s](%s)\n", i, title, base_url + notice_url)
	})
}

func monitornotices(msg *tgba.MessageConfig) {
	f, err := os.Open("scratchpad/diff.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}
	base_url := "https://rlacollege.edu.in/"
	// use the goquery document...
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		notice_url, _ := s.Attr("href")

		parts := strings.Split(notice_url, "/")
		// ) is wrongfully parsed as markdown syntax, resulting in broken links in msg
		parts[2] = url.QueryEscape(parts[2])
		notice_url = strings.Join(parts[:], "/")

		title := s.Text()
		msg.Text += fmt.Sprintf("%d: [%s](%s)\n", i, title, base_url + notice_url)
	})
	
}

func main() {
	bot, err := tgba.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgba.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgba.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyToMessageID = update.Message.MessageID
			switch update.Message.Command() {
				case "help": msg.Text = `Type one of the following commands:
				/help - display this message
				/sayhi - make me say hi
				/status - check if I'm ok
				/website - send the college website link
				/notices - send notices from the website`

				case "sayhi": msg.Text = "Hi! :)"

				case "status": msg.Text = "I'm ok, thanks for concern."

				case "website":
					msg.ParseMode = "Markdown"
					msg.Text = "[RLA Website](https://rlacollege.edu.in/)"

				case "notices": getNotices(&msg)

				case "newnotice":
					msg.ParseMode = "Markdown"
					monitornotices(&msg)
					//time.Sleep(20 * time.Second)

				default: msg.Text = "I don't know that command"
			}

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

