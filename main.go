package main

import (
	"Goyou/Goyou"
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
	botToken := os.Getenv("bott") //token["botToken"]

	b, err := tb.NewBot(tb.Settings{Token: botToken,
		ParseMode: tb.ModeMarkdownV2, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})

	if err != nil {
		log.Fatalln(err)
		return
	}
	var query []Goyou.Details

	b.Handle("/help", func(m *tb.Message) {
		b.Reply(m, "Nothing to show")
	})
	b.Handle("/start", func(m *tb.Message) {
		b.Reply(m, fmt.Sprintf("Hello *%s*, welcome to *%s*\n 1\\. /start : To start the bot\n 2\\. /help : Get infomation about bot\\.", m.Sender.Username, b.Me.FirstName))

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
			query = Goyou.Query(url)
			buttons := make([]tb.Row, len(query))
			menu := &tb.ReplyMarkup{}

			for i, video := range query {
				var btn tb.Row

				if i%2 == 0 {
					if len(query) != i+1 {
						var details = query[i+1]
						btn = menu.Row(menu.Data(
							fmt.Sprintf("%s, %s", video.Quality, video.Size),
							video.Quality+video.Size, strconv.Itoa(i)),
							menu.Data(
								fmt.Sprintf("%s, %s", details.Quality, details.Size),
								details.Quality+details.Size, strconv.Itoa(i+1)))

					}
				} else if len(query)%2 != 0 && i == (len(query)-1) {
					btn = menu.Row(menu.Data(
						fmt.Sprintf("%s, %s", video.Quality, video.Size),
						video.Quality+video.Size, strconv.Itoa(i)))
				}
				buttons = append(buttons, btn)

			}

			//menu := &tb.ReplyMarkup{}
			//btn := menu.Data("HEE", "154", url)
			//btn2 := menu.Data("Nawa", "kkk", url)

			//tb.Btn{}
			//menu.Inline(menu.Row(btn), menu.Row(btn2))
			menu.Inline(buttons...)
			photo := tb.Photo{File: tb.FromURL(query[1].Thumb.URL), Width: 400, Height: 400}
			_, err2 := b.Reply(m, &photo, &tb.ReplyMarkup{InlineKeyboard: menu.InlineKeyboard})
			if err2 != nil {
				fmt.Println(err2)
				return
			}
			return
		} else if len(strings.TrimSpace(m.Text)) == 3 {
			b.Reply(m, "YouTube link not specified\\.\n*Usage:*\n\\-\\ /yt video link\\.")
			return
		}
		b.Reply(m, "*Not a youtube link*")
	})
	b.Handle(tb.OnCallback, func(c *tb.Callback) {
		if len(query) != 0 {
			b.Respond(c, &tb.CallbackResponse{})
			b.Delete(c.Message)
			mId, err := b.Reply(c.Message.ReplyTo, "*Downloading*")
			if err != nil {
				log.Println(err)
			}

			item, _ := strconv.Atoi(c.Data)
			filename := Goyou.Download(query[item])
			mId2, err1 := b.Edit(mId, "*Uploading*")
			if err1 != nil {
				log.Println(err1.Error())
			}
			photo := &tb.Photo{File: tb.FromDisk(filename)}
			filed := tb.Video{
				File:      tb.FromDisk(filename),
				Thumbnail: photo,
				FileName:  filename,
				Caption:   query[item].Title,
				MIME:      query[item].Stream.MimeType,
			}
			err2 := b.Notify(c.Message.Chat, tb.UploadingVideo)
			if err2 != nil {
				log.Println(err2.Error())
				return
			}
			b.Delete(mId2)
			_, err3 := b.Reply(c.Message.ReplyTo, &filed, tb.ModeHTML)
			if err3 != nil {
				log.Panic(err3.Error())
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
