package photosafe

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

type PhotoTemp struct {
	File            string
	Caption         string
	CaptionEntities []tgbotapi.MessageEntity
}

func PhotoSafe(adminId []int64) {
	bot, err := tgbotapi.NewBotAPI("7071645488:AAFbkH6wGo7OUPRpJ6cHs5PWrK1ryN4Mkrk")
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		flag := false
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
			panic(err)
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
}

func ChangePhoto(update tgbotapi.Update) {
	file, err := os.Open("inf.txt")
	if err != nil {
		panic(err)
	}

	var photoTemp PhotoTemp
	err = json.NewDecoder(file).Decode(&photoTemp)
	if err != nil {
		panic(err)
	}
	file.Close()

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
	if err != nil {
		panic(err)
	}

	var photoTemp PhotoTemp
	err = json.NewDecoder(file).Decode(&photoTemp)
	if err != nil {
		panic(err)
	}
	file.Close()

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
