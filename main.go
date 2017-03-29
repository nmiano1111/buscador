package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"time"
)

type User struct {
	ID               int    `json:"id"`
	LastLoginTime    string `json:"last_login_time"` //time.Time `json:"last_login_time"`
	Name             string `json:"nickname"`
	Intro            string `json:"intro"`
	LearningLanguage string `json:"learning_language"`
	LivingIn         string `json:living_country_id`
	From             string `json:"origin_country_id"`
}

type UserWrapper struct {
	Users []User
	Done  bool
	Error error
}

type Page struct {
	Data []User `json:"data"`
	Meta struct {
		HasNext bool `json:"has_next"`
	} `json:"meta"`
}

const baseURL = "https://www.italki.com"

func main() {
	//i_token=TWpnek5qWTNNZz09fDE0OTAxOTMyMDR8MjFkODY3OTk5ODQ2NmFjMWYzMjg2MGRhNzFiNWY1MDY0ZTJhZDU1OA%3D%3D
	//
	users, err := getAllUsers()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(users))
}

func getAllUsers() ([]User, error) {
	userChan := make(chan UserWrapper)
	go getUsers(1, userChan)

	users := []User{}
	for user := range userChan {
		if user.Error != nil {
			return users, nil
		}
		users = append(users, user.Users...)
		if user.Done {
			break
		}
	}

	return users, nil
}

/*
 1. start at page 1.
 2. get users
 3. check if 'has_next'
 4. if 'has_next', recursively get next page on new goroutine
*/
func getUsers(index int, userChan chan UserWrapper) {
	partnerURL := fmt.Sprintf("%s/api/partner?_r=1490718821186&city=&country=CL&gender=1&hl=en-US&is_native=1&learn=english&page=%d&speak=spanish", baseURL, index)
	resp, err := http.Get(partnerURL)
	if err != nil {
		userChan <- UserWrapper{Error: err}
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		userChan <- UserWrapper{Error: err}
	}

	var page Page
	err = json.Unmarshal(bytes, &page)
	if err != nil {
		userChan <- UserWrapper{Error: err}
	}
	fmt.Print(".")
	userChan <- UserWrapper{Users: page.Data}
	if page.Meta.HasNext {
		go getUsers(index+1, userChan)
	} else {
		userChan <- UserWrapper{Done: true}
	}
}
