package Conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type BotCofig struct {
	ClientSecret string
	ClientID     string
	TgApi        string
	LogSt        string
	ChatId       int64
	AdminId      []int64
}

func Config() *BotCofig {
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
	return &conf
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
	fmt.Print("Введите login стримера, которога желаете отслеживать: ")
	_, err = fmt.Scan(&conf.LogSt)
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
