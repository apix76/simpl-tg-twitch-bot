package photosafe

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"time"
)

type TempMess struct {
	File            string
	Caption         string
	Entities        []tgbotapi.MessageEntity
	CaptionEntities []tgbotapi.MessageEntity
}

func WaitMess(adminId []int64, TgApi string) {
	for {
		err := PhotoSafe(adminId, TgApi)
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
		}
	}
}

func PhotoSafe(adminId []int64, TgApi string) error {
	bot, err := tgbotapi.NewBotAPI(TgApi)
	if err != nil {
		return err
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updateConfig.AllowedUpdates = []string{"message"}
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		flag := false
		if update.Message == nil {
			continue
		}
		fmt.Printf("%+v\n", update.Message)
		for _, v := range adminId {
			if v == update.Message.From.ID {
				flag = true
			}
			if update.Message.Chat.ID != update.Message.From.ID {
				continue
			}
		}
		if flag != true {
			continue
		}
		if update.Message == nil {
			continue
		}
		mes := tgbotapi.NewMessage(update.Message.Chat.ID, "Принял")
		_, err := bot.Send(mes)
		if err != nil {
			return err
		}
		if update.Message.Photo != nil && update.Message.Caption == "" {
			ChangePhoto(update)
		}
		if update.Message.Photo == nil && update.Message.Text != "" {
			ChangeText(update)
		}
		if update.Message.Photo != nil && update.Message.Caption != "" {
			ChangeAll(update)
		}
	}
	return nil
}

func ChangePhoto(update tgbotapi.Update) {
	file, err := os.Open("inf.txt")
	var photoTemp TempMess
	if err == nil {
		err = json.NewDecoder(file).Decode(&photoTemp)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	photo := tgbotapi.PhotoConfig{
		Caption:         photoTemp.Caption,
		CaptionEntities: photoTemp.CaptionEntities,
		BaseFile: tgbotapi.BaseFile{
			File: tgbotapi.FileID(update.Message.Photo[0].FileID),
		},
	}

	file, err = os.Create("inf.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	photobyte, err := json.Marshal(photo)
	if err != nil {
		panic(err)
	}
	file.Write(photobyte)

}

func ChangeText(update tgbotapi.Update) {
	file, err := os.Open("inf.txt")
	var photoTemp TempMess
	if err == nil {
		err = json.NewDecoder(file).Decode(&photoTemp)
		if err != nil {
			panic(err)
		}
		file.Close()
	}
	photo := tgbotapi.PhotoConfig{
		Caption:         update.Message.Text,
		CaptionEntities: update.Message.Entities,
		BaseFile: tgbotapi.BaseFile{
			File: tgbotapi.FileID(photoTemp.File),
		},
	}

	file, err = os.Create("inf.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	photobyte, err := json.Marshal(photo)
	if err != nil {
		panic(err)
	}
	file.Write(photobyte)
}

func ChangeAll(update tgbotapi.Update) {
	photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileID(update.Message.Photo[0].FileID))
	photo.Caption = update.Message.Caption
	photo.CaptionEntities = update.Message.CaptionEntities

	file, err := os.Create("inf.txt")
	if err != nil {
		panic(err)
	}
	photoBute, err := json.Marshal(photo)
	if err != nil {
		panic(err)
	}
	file.Write(photoBute)
	file.Close()
}
