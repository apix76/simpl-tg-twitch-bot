package main

import (
	photosafe "awesomeProject3/Bot_testg"
	"awesomeProject3/Conf"
	TwitchAccessTocen "awesomeProject3/twitch"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicklaw5/helix"
	"log"
	"os"
	"time"
)

func main() {
	conf := Conf.Config()
	go photosafe.WaitMess(conf.AdminId, conf.TgApi)
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
			UserLogins: []string{conf.LogSt},
		})
		if err != nil {
			log.Printf("Resp get errors! : %v\n", err)
			client.SetAppAccessToken(TwitchAccessTocen.AccessToken(conf.ClientID, conf.ClientSecret))
			time.Sleep(5 * time.Second)
			continue
		}
		if RedFlag == false && len(resp.Data.Streams) != 0 {
			streamerFlag := false
			for _, v := range resp.Data.Streams {
				if v.UserLogin == conf.LogSt {
					streamerFlag = true
				}
			}
			if !streamerFlag {
				continue
			}

			fmt.Printf("%+v\n", resp)
			file, err := os.Open("inf.txt")
			var Tempmess photosafe.TempMess
			if err != nil {
				fmt.Printf("Error open file, check available file: %v", err)
			} else {
				err = json.NewDecoder(file).Decode(&Tempmess)
				if err != nil {
					panic(err)
				}
			}
			file.Close()
			if Tempmess.File != "" {
				MessInf, err = Photo(bot, Tempmess, *conf, resp.Data.Streams[0].GameName)
				if err != nil {
					log.Printf("MessInf errors (Photo)! : %v\n", err)
					continue
				}
			} else {
				MessInf, err = Mess(bot, Tempmess, *conf, resp.Data.Streams[0].GameName)
				if err != nil {
					log.Printf("MessInf errors (Mess)! : %v\n", err)
					continue
				}
			}

			RedFlag = true
		}
		if RedFlag == true && len(resp.Data.Streams) == 0 {
			RedFlag = false
			del := tgbotapi.NewDeleteMessage(conf.ChatId, MessInf.MessageID)
			if _, err := bot.Request(del); err != nil {
				log.Printf("Error delet message. May be he is deleted! : %v\n", err)
				continue
			}
		}

		time.Sleep(30 * time.Second)
	}
}

func Photo(bot *tgbotapi.BotAPI, TempMess photosafe.TempMess, con Conf.BotCofig, GameName string) (tgbotapi.Message, error) {
	photo := tgbotapi.PhotoConfig{
		Caption:         TempMess.Caption,
		CaptionEntities: TempMess.CaptionEntities,
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID: con.ChatId,
			},
			File: tgbotapi.FileID(TempMess.File),
		},
	}

	photo.Caption = fmt.Sprintf("%v\n\nКатегория: %v\n#стрим #twitch", photo.Caption, GameName)

	InlineButten := tgbotapi.NewInlineKeyboardButtonURL("Смотреть стрим!", fmt.Sprintf("https://www.twitch.tv/%v", con.LogSt))
	InlineButtenRow := tgbotapi.NewInlineKeyboardRow(InlineButten)
	InlineKeyBordeMarkup := tgbotapi.NewInlineKeyboardMarkup(InlineButtenRow)
	photo.BaseChat.ReplyMarkup = InlineKeyBordeMarkup

	idmess, err := bot.Send(photo)
	return idmess, err
}

func Mess(bot *tgbotapi.BotAPI, TempMess photosafe.TempMess, con Conf.BotCofig, GameName string) (tgbotapi.Message, error) {
	mess := tgbotapi.MessageConfig{
		Text:     TempMess.Caption,
		Entities: TempMess.CaptionEntities,
		BaseChat: tgbotapi.BaseChat{
			ChatID: con.ChatId,
		},
	}

	mess.Text = fmt.Sprintf("%v\n\nКатегория: %v\n#стрим #twitch", mess.Text, GameName)

	InlineButten := tgbotapi.NewInlineKeyboardButtonURL("Смотреть стрим!", fmt.Sprintf("https://www.twitch.tv/%v", con.LogSt))
	InlineButtenRow := tgbotapi.NewInlineKeyboardRow(InlineButten)
	InlineKeyBordeMarkup := tgbotapi.NewInlineKeyboardMarkup(InlineButtenRow)
	mess.BaseChat.ReplyMarkup = InlineKeyBordeMarkup

	idmess, err := bot.Send(mess)
	return idmess, err
}
