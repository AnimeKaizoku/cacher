package main

import (
	"fmt"
	"time"

	"github.com/AnimeKaizoku/cacher"
)

// This program is an example of Type 2 SingleKey System
//
// It will be a int-struct mapping, where int key would
// be user id of our hypothetical user for some storage service,
// and struct value would contain user's name (type string), dob
// (type struct{date int, month int, year int}).
//
// We suggest you to run this program, edit, experiment with it
// to know its working better.
// Suggested experiment:
//    Try to fetch user id 100 within 10 seconds of starting the program
//    then search the same user id again after 5 seconds and wait for
//    another 5 seconds, and then try any other id except the one you
//    tried earlier and then immediately try user id 100.
// Observation:
// 	  You'll observe that user id 100 is still present in cache while
// 	  others have expired, this happens because of revaluation mode.
// 	  More frequently a key is used, later will be its expiration.

type User struct {
	Name string
	Dob  date
}

type date struct{ date, month, year int }

// We'll set expiration (timeToLive) to 10 seconds
// And our cleaner window to run in every 2 minutes
// We will also keep the revaluation mode on so that
// the if keys are being used frequently then they
// shouldn't expire soon.
var cache = cacher.NewCacher[int, User](&cacher.NewCacherOpts{
	TimeToLive:    time.Second * 10,
	CleanInterval: time.Minute * 2,
	Revaluate:     true,
})

// This map contains our sample date
var data = map[int]User{
	100: {"John", date{15, 05, 1976}},
	101: {"Stuart", date{20, 01, 2000}},
	102: {"Susan", date{29, 02, 2016}},
}

func main() {
	// This will add all the hypothetical case to our cache instance
	for userId, userInfo := range data {
		cache.Set(userId, userInfo)
	}

	// Some fancy lines
	fmt.Println("KlingonStorage: A hypothetical storage service")
	fmt.Println("I am a program to help you get details of a user from his userid")
	for {
		ask()
	}
}

func ask() {
	fmt.Print("Enter User ID: ")

	var userId int
	// Scanln will scan the input and populate our userId variable
	fmt.Scanln(&userId)

	// value variable will contain the title of chat while ok
	// variable will tell us whether it was found unexpired in cache not.
	value, ok := cache.Get(userId)
	if !ok {
		fmt.Println("Provided user id has either expired or was not found in the cache!")
		return
	}
	fmt.Println("Details of this user are:", value, "\n ")
}
