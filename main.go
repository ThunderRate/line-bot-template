package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/spf13/viper"
)

var bot *linebot.Client

func exit(err error) {
	var status int
	if err != nil {
		log.Printf("Exit with error: %s\n", err.Error())
		status = 2
	}
	os.Exit(status)
}

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		exit(err)
	}
}

func main() {
	var err error
	if bot, err = linebot.New(viper.GetString("ChannelSecret"), viper.GetString("ChannelAccessToken")); err != nil {
		exit(err)
	}
	addr := net.JoinHostPort(viper.GetString("Host"), viper.GetString("Port"))
	http.HandleFunc("/callback", callbackHandler)
	if err = http.ListenAndServe(addr, nil); err != nil {
		exit(err)
	}
	exit(nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(200)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
