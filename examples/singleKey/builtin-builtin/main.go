package main

import (
	"fmt"
	"time"

	"github.com/AnimeKaizoku/cacher"
)

// This program is an example of Type 1 SingleKey System
//
// It will be a int64-string mapping, where int64 key would
// be chat id of our hypothetical chatting app, and string
// value would be the title of that chat.

// We'll set expiration (timeToLive) to 10 seconds
// And our cleaner window to run in every 2 minutes
// We will also keep the revaluation mode on so that
// the if keys are being used frequently then they
// shouldn't expire soon.
var cache = cacher.NewCacher[int64, string](&cacher.NewCacherOpts{
	TimeToLive:    time.Second * 10,
	CleanInterval: time.Minute * 2,
	Revaluate:     true,
})

// Data we will be caching for our hypothetical app
var data = map[int64]string{
	100100228211: "Chatter Support",
	100228330210: "Chatter OT",
	100228337922: "Dear Diary",
	100883399228: "Random Chat",
}

func main() {
	// This will add all the hypothetical case to our cache instance
	for chatId, title := range data {
		cache.Set(chatId, title)
	}

	// Some fancy lines
	fmt.Println("ChatterGram: A hypothetical chatting app")
	fmt.Println("I am a program to help you get title from a chatid")
	for {
		ask()
	}
}

func ask() {
	fmt.Print("Enter Chat ID: ")

	var chatId int64
	// Scanln will scan the input and populate our chatId variable
	fmt.Scanln(&chatId)

	// value variable will contain the title of chat while ok
	// variable will tell us whether it was found unexpired in cache not.
	value, ok := cache.Get(chatId)
	if !ok {
		fmt.Println("Provided chat id has either expired or was not found in the cache!")
		return
	}
	fmt.Println("Title for this chat is:", value)
}
