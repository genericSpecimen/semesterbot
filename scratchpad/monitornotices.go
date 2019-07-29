package main

import (
	"fmt"
	"log"
	"time"
	"os/exec"
	"io/ioutil"
	"os"
	"strings"
	"net/http"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func monitor(noticeupdate chan string) {
	log.Printf("Started monitoring notices..")
	for {
		if err := exec.Command("/bin/sh", "./makediff.sh").Run(); err != nil {
			log.Printf("Failed to execute makediff.sh")
		}
		log.Printf("Successfully executed makediff.sh")

		data, err := ioutil.ReadFile("diff.html")
		if err != nil {
			log.Panic(err)
		}

		da := string(data)

		/*
		fmt.Print(da)
		fmt.Println(len(da))
		*/

		if(len(da) > 0) {
			log.Printf("New notices were found.")
			noticeupdate <- da

			cmd := "cp"
			args := []string{"new.php", "old.php"}
			if err := exec.Command(cmd, args...).Run(); err != nil {
				log.Printf("Failed: cp new.php old.php")
			}
			
		} else {
			log.Printf("No new notices were found.")
		}

		time.Sleep(10 * time.Second)
	}
	log.Printf("Stopped monitoring notices..")
}

func getAllNotices(msg *tgba.MessageConfig) {
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

	noticeupdate := make(chan string)
	go monitor(noticeupdate)

	for {
		select {
			case notices := <-noticeupdate: {
				log.Printf("New notice received: %s", notices)
				msg := tgba.NewMessage(-288929399, "")
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(notices))
				if err != nil {
					log.Panic(err)
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

				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}

			case update := <-updates: {
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
		
						case "notices": getAllNotices(&msg)
						
						default: msg.Text = "I don't know that command"
					}
		
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
			}
		}
	}
}
