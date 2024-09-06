package main

import (
	photosafe "awesomeProject3/Bot_testg"
	TwitchAccessTocen "awesomeProject3/twitch"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicklaw5/helix"
	"os"
	"time"
)

type BotCofig struct {
	ClientSecret string
	ClientID     string
	TgApi        string
	ChatId       int64
	AdminId      []int64
}

func main() {
	var conf BotCofig
	fileCon, err := os.Open("config.cfg")
	if err != nil {
		fileCon, err = os.Create("config.cfg")
		if err != nil {
			panic(err)
		}
		NewConfig(&conf)
		confbyte, err := json.Marshal(conf)
		if err != nil {
			panic(err)
		}
		_, err = fileCon.Write(confbyte)
		if err != nil {
			panic(err)
		}
		fileCon.Close()
	} else {
		err = json.NewDecoder(fileCon).Decode(&conf)
		if err != nil {
			fileCon.Close()
			fileCon, err = os.Create("config.cfg")
			NewConfig(&conf)
			confbyte, err := json.Marshal(conf)
			if err != nil {
				panic(err)
			}
			_, err = fileCon.Write(confbyte)
			if err != nil {
				panic(err)
			}
			fileCon.Close()
		}
	}

	go photosafe.PhotoSafe(conf.AdminId)
	bot, err := tgbotapi.NewBotAPI(conf.TgApi)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	bot.Debug = true
	RedFlag := false

	var MessInf tgbotapi.Message

	AccessToken := TwitchAccessTocen.AccessToken(conf.ClientID, conf.ClientSecret)
	client, err := helix.NewClient(&helix.Options{
		ClientID:       conf.ClientID,
		AppAccessToken: AccessToken,
	})
	for {
		resp, err := client.GetStreams(&helix.StreamsParams{
			UserLogins: []string{"daratama_"},
		})
		if err != nil {
			// handle error
		}
		if RedFlag == false && len(resp.Data.Streams) != 0 {
			fmt.Printf("%+v\n", resp)
			file, err := os.Open("inf.txt")
			if err != nil {
				panic(err)
			}
			var photoTemp photosafe.PhotoTemp
			err = json.NewDecoder(file).Decode(&photoTemp)
			if err != nil {
				panic(err)
			}
			photo := tgbotapi.PhotoConfig{
				Caption:         photoTemp.Caption,
				CaptionEntities: photoTemp.CaptionEntities,
				BaseFile: tgbotapi.BaseFile{
					BaseChat: tgbotapi.BaseChat{
						ChatID: conf.ChatId,
					},
					File: tgbotapi.FileID(photoTemp.File),
				},
			}
			photo.Caption = fmt.Sprintf("%v\n\nКатегория: %v\n#стрим #twitch", photo.Caption, resp.Data.Streams[0].GameName)

			InlineButten := tgbotapi.NewInlineKeyboardButtonURL("Смотреть стрим!", "https://www.twitch.tv/daratama_")
			InlineButtenRow := tgbotapi.NewInlineKeyboardRow(InlineButten)
			InlineKeyBordeMarkup := tgbotapi.NewInlineKeyboardMarkup(InlineButtenRow)
			photo.BaseChat.ReplyMarkup = InlineKeyBordeMarkup

			MessInf, err = bot.Send(photo)
			if err != nil {
				panic(err)
			}
			RedFlag = true
		}
		if RedFlag == true && len(resp.Data.Streams) == 0 {
			del := tgbotapi.NewDeleteMessage(conf.ChatId, MessInf.MessageID)
			if _, err := bot.Request(del); err != nil {
				// Note that panics are a bad way to handle errors. Telegram can
				// have service outages or network errors, you should retry sending
				// messages or more gracefully handle failures.
				panic(err)
			}
			RedFlag = false
		}

		time.Sleep(30 * time.Second)
	}
}

func NewConfig(conf *BotCofig) {
	fmt.Print("Введите client secret приложения twitch: ")
	_, err := fmt.Scan(&conf.ClientSecret)
	if err != nil {
		panic(err)
	}
	fmt.Print("Введите client id приложения twitch: ")
	_, err = fmt.Scan(&conf.ClientID)
	if err != nil {
		panic(err)
	}
	fmt.Print("Введите chatid из telegram: ")
	_, err = fmt.Scan(&conf.ChatId)
	if err != nil {
		panic(err)
	}
	fmt.Print("Введите api token от BotFather из telegram: ")
	_, err = fmt.Scan(&conf.TgApi)
	if err != nil {
		panic(err)
	}
	fmt.Println("Далее заполните id пользователей Telegram которые могут модерировать бота.\nДля выхода введите 0")
	for {
		var id int64
		fmt.Print("id: ")
		_, err := fmt.Scan(&id)
		if err != nil {
			panic(err)
		}
		if id == 0 {
			break
		}
		conf.AdminId = append(conf.AdminId, id)
	}

}
