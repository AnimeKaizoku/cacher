package main

import (
	"fmt"
	"strings"

	"github.com/AnimeKaizoku/cacher"
)

// We'll be using PolyKeyer hence we need to choose string as the data type for Key.
// In this example, we won't be expiring keys and hence we don't pass them to opts.
var cache = cacher.NewCacher[string, string](nil)

// Following code initiates a PolyKeyer instance.
// We've chosen "chat" as the prefix key (also called primary key) and 2 as our secondary
// parameters.
// Note that it'll panic in case you pass different number of secondary keys to the keyer.New
// method than the chosen one while initiation.
var keyer = cacher.NewPolyKeyer("chat", 2)

// This function will call Get method on input key and print it if found.
func get(key string) {
	title, ok := cache.Get(key)
	if !ok {
		fmt.Printf("Key '%s' not found in the valid cache!\n", key)
		return
	}
	fmt.Println("Title of that chat is:", title)
}

// This function will fetch all the keys in our current cacher and print them.
func getAll() {
	allKeys := cache.GetAll()
	fmt.Println("All keys:")
	for _, x := range allKeys {
		fmt.Println("-", x)
	}
}

func main() {
	// We create a new key with exactly 2 arguments (we chose while initialising polykeyer)
	// Let's assume that our 1st argument will be chat id and 2nd argument can either be
	// "public" or "private".
	var key = keyer.New("101", "private")
	// Here we set the value for our key we created above and let our value be the title
	// of chat.
	cache.Set(key, "King's Chat")

	// Uncomment the following line of code to see how does a poly key look like actually:
	// fmt.Println("This is how a polykey look like:", key)

	// Let's retrieve our key
	// We'll use the get function we wrote earlier in this example for it!
	get(key)

	// We can know the number of keys which are present in our current
	// Cacher instance using the cacher.NumKeys() method.
	fmt.Println("Number of keys in current cacher:", cache.NumKeys())

	// Resetting the cache will delete all key-value mapping present in
	// it currently.
	cache.Reset()

	// Let's try to get our key now, it should not be found as we recently
	// reset our cache.
	get(key)

	// Let's add the key again and a few more for further tutorial.
	cache.Set(key, "King's Chat")
	cache.Set(keyer.New("102", "public"), "Bob's chitchat group")
	cache.Set(keyer.New("103", "private"), "Cacher Test Chat")

	key1 := keyer.New("104", "public")
	cache.Set(key1, "Github Public Chat")

	cache.Set(keyer.New("105", "private"), "King Hero")

	// Let's print the number of keys now:
	fmt.Println("Number of keys:", cache.NumKeys())

	// Now we will learn how to make a segrigator function which we can
	// use for doing ___Some calls which are conditional in nature.
	// Let out condition be: return true for all values which start
	// with "King"
	exampleSegrigator := func(v string) bool {
		// We'll use HasPrefix function from strings package determining
		// prefix.
		return strings.HasPrefix(v, "King")
	}

	// We can also get all or some particular keys in an array at once.
	// Following is the code on how to get some keys.
	values := cache.GetSome(exampleSegrigator)

	// Let's print our conditional values
	fmt.Println(`Here are chats which start with the phrase "King"`)
	for _, x := range values {
		fmt.Println("-", x)
	}

	// Now we'll learn how to delete a specific key from our current cacher
	// Let's suppose we want to delete key with chatid 104 and "public"
	// secondary keys, we'll use the key1 variable we defined earlier for
	// the purpose.
	cache.Delete(key1)

	// Let's print all the keys currently present in the cacher
	// We'll use getAll function we created earlier in this example for it.
	getAll()

	// We've learnt a lot till yet, it's the time we introduce ourselves to
	// the DeleteSome method, similar to GetSome, it deletes keys on the basis
	// of segrigations, we'll just use the previous segrigator function and delete
	// all keys which start with the phrase "King"
	cache.DeleteSome(exampleSegrigator)

	// Now we will again use getAll to see the difference finally
	getAll()
}

// Congratulations! You just learnt most of the usable methods of this library!
// We request you to check out Single Key System examples to learn how to make
// a cacher which will expire keys and that too with the revaluation mode.
