package main

import (
	"fmt"
	"time"

	"github.com/AnimeKaizoku/cacher"
)

// This program is an example of Type 3 SingleKey System
//
// It will be a (custom int)-([]struct)  mapping, where
// custom int will be rank level of a person in an agency
// from (1, 2, 3) and value array will contain struct
// consisting the user's name and dob.

type rank int

func (p rank) Stringify() string {
	switch p {
	case 1:
		return "Admiral"
	case 2:
		return "Captain"
	case 3:
		return "Ensign"
	default:
		return "invalid"
	}
}

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
var cache = cacher.NewCacher[rank, []User](&cacher.NewCacherOpts{
	TimeToLive:    time.Second * 10,
	CleanInterval: time.Minute * 2,
	Revaluate:     true,
})

// This map contains our sample date
var data = map[rank][]User{
	1: {
		{"JL Picard", date{15, 05, 1976}},
		{"Tom Riker", date{21, 03, 1988}},
	},
	2: {
		{"Michael", date{15, 05, 1998}},
		{"James Kirk", date{29, 03, 1995}},
	},
	3: {
		{"John Doe", date{15, 05, 2004}},
		{"Harry Potter", date{12, 04, 2006}},
	},
}

func main() {
	// This will add all the hypothetical case to our cache instance
	for userId, userInfo := range data {
		cache.Set(userId, userInfo)
	}

	// Some fancy lines
	fmt.Println("GopherSec: A hypothetical agency")
	fmt.Println("I am a program to help you get details of a user from his userid")
	for {
		ask()
	}
}

func ask() {
	fmt.Print("Enter User ID: ")

	var rank rank
	// Scanln will scan the input and populate our userId variable
	fmt.Scanln(&rank)

	// value variable will contain the title of chat while ok
	// variable will tell us whether it was found unexpired in cache not.
	value, ok := cache.Get(rank)
	if !ok {
		fmt.Println("Info of this rank has either expired or was not found in the cache!")
		return
	}
	fmt.Println("Cached Info:")
	fmt.Println("\nRank:", rank.Stringify())
	fmt.Println("Name : DOB")
	for _, x := range value {
		fmt.Println(x.Name, ":", x.Dob)
	}
	fmt.Println("Search Complete!\n ")
}
