package main

import (
	"fmt"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//token, _ := godotenv.Read()
	if errdot := godotenv.Load(); errdot != nil {
		log.Fatal("Error loading .env file")
	}
	botToken := os.Getenv("botToken") //token["botToken"]

	b, err := tb.NewBot(tb.Settings{Token: botToken,
		ParseMode: tb.ModeMarkdownV2, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})

	if err != nil {
		log.Fatalln(err)
		return
	}
	var query []Details

	b.Handle("/start", func(m *tb.Message) {
		b.Reply(m, fmt.Sprintf("*Welcome* %s to *Zaza*\n 1\\. /start : To start the bot\n 2\\. /help : Get infomation about bot\\.", m.Sender.Username))

	})
	b.Handle(tb.OnText, func(m *tb.Message) {
		if strings.ToLower(m.Text) == "kill" {
			b.Reply(m, "*Bot Killed*")
			b.Stop()
			return
		}
		b.Reply(m, m.Text)

	})
	b.Handle("/yt", func(m *tb.Message) {
		if strings.Contains(m.Text, "youtube") || strings.Contains(m.Text, "youtu") {
			url := m.Text[4:]
			query = Query(url)
			buttons := make([]tb.Row, len(query))
			menu := &tb.ReplyMarkup{}

			for i, video := range query {
				var btn tb.Row

				if i%2 == 0 {
					if len(query) != i+1 {
						var details = query[i+1]
						btn = menu.Row(menu.Data(
							fmt.Sprintf("%s, %s", video.quality, video.size),
							video.quality+video.size, strconv.Itoa(i)),
							menu.Data(
								fmt.Sprintf("%s, %s", details.quality, details.size),
								details.quality+details.size, strconv.Itoa(i+1)))

					}
				} else if len(query)%2 != 0 && i == (len(query)-1) {
					btn = menu.Row(menu.Data(
						fmt.Sprintf("%s, %s", video.quality, video.size),
						video.quality+video.size, strconv.Itoa(i)))
				}
				buttons = append(buttons, btn)

			}

			//menu := &tb.ReplyMarkup{}
			//btn := menu.Data("HEE", "154", url)
			//btn2 := menu.Data("Nawa", "kkk", url)

			//tb.Btn{}
			//menu.Inline(menu.Row(btn), menu.Row(btn2))
			menu.Inline(buttons...)
			photo := tb.Photo{File: tb.FromURL(query[1].thumb.URL), Width: 400, Height: 400}
			_, err2 := b.Reply(m, &photo, &tb.ReplyMarkup{InlineKeyboard: menu.InlineKeyboard})
			if err2 != nil {
				fmt.Println(err2)
				return
			}
			return
		}
		b.Reply(m, "*Not a youtube link*")
	})
	b.Handle(tb.OnCallback, func(c *tb.Callback) {
		if len(query) != 0 {
			b.Respond(c, &tb.CallbackResponse{})
			item, _ := strconv.Atoi(c.Data)
			filename := Download(query[item])
			b.Reply(c.Message, "Downloaded", tb.ModeHTML)
			photo := &tb.Photo{File: tb.FromDisk(filename)}
			filed := tb.Video{
				File:      tb.FromDisk(filename),
				Thumbnail: photo,
				FileName:  filename,
				Caption:   query[item].title,
				MIME:      query[item].stream.MimeType,
			}
			err1 := b.Notify(c.Message.Chat, tb.UploadingVideo)
			if err1 != nil {
				log.Println(err1.Error())
				return
			}

			_, err2 := b.Send(c.Sender, &filed, tb.ModeHTML)
			if err2 != nil {
				log.Panic(err2.Error())
				return
			}

		}
		//b.Respond(c, &tb.CallbackResponse{})
		//b.Reply(c.Message, "This message is old. resend")
	})
	log.Println("Bot Started")
	b.Start()

}

/*inlineKeys := [][]tb.InlineButton{{tb.InlineButton{ //https://github.com/tucnak/telebot/blob/v2.5.0/callbacks.go#L65
	Unique: "foo_btnp",
	Text:   "foo",
	Data:   "foo_btn",
}, tb.InlineButton{
	Unique: "bar_btnp",
	Text:   "btn",
	Data:   "bar_btn",
},
}}
*/
