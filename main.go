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
	count := 1
	for {
		if count <= 59 {
			go getUsers(count, userChan)
			count = count + 1
		} else {
			break
		}
	}

	users := []User{}
	for user := range userChan {
		if user.Error != nil {
			return users, nil
		}
		if user.Done {
			break
		}
		fmt.Println("----")
		fmt.Println(len(users))
		fmt.Println("----")
		users = append(users, user.Users...)
	}

	return users, nil
}

func getUsers(page int, userChan chan UserWrapper) {
	partnerURL := fmt.Sprintf("%s/api/partner?_r=1490718821186&city=&country=CL&gender=1&hl=en-US&is_native=1&learn=english&page=%d&speak=spanish", baseURL, page)
	resp, err := http.Get(partnerURL)
	if err != nil {
		userChan <- UserWrapper{Error: err}
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		userChan <- UserWrapper{Error: err}
	}

	type wrapper struct {
		Data []User `json:"data"`
	}

	var wrappedUsers wrapper
	err = json.Unmarshal(bytes, &wrappedUsers)
	if err != nil {
		userChan <- UserWrapper{Error: err}
	}
	userChan <- UserWrapper{Users: wrappedUsers.Data}
}
