// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The simple command demonstrates a simple functionality which
// prompts the user for a GitHub username and lists all the public
// organization memberships of the specified username.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v40/github"

	"golang.org/x/oauth2"
)

func main() {
	//fmt.Print("Enter User Location : ")
	//var input string
	// fmt.Scanln(&input)
	// fmt.Print(input)
	//input = "Italy"

	ctx := context.Background()

	var token string = os.Getenv("PAT")
	fmt.Println(token)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}
	// get all pages of results
	var allUsers []*github.User
	var userData []string
	// https://docs.github.com/en/rest/reference/search#search-users
	//query := "user:crixo"
	query := "language:c# location:italy"

	for {
		userSearch, resp, err := client.Search.Users(ctx, query, opts)
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("Page %d of %d users.", opts.Page, opts.PerPage))
		//fmt.Println(len(userSearch.Users))

		for _, user := range userSearch.Users {
			allUsers = append(allUsers, user)
			fullUser, _ := fetchUsers(client, *user.ID)

			var userProps []string

			if fullUser.Login != nil {
				userProps = append(userProps, fmt.Sprintf("Login: %s", *fullUser.Login))
			}

			if fullUser.Name != nil {
				userProps = append(userProps, fmt.Sprintf("Name: %s", *fullUser.Name))
			}
			if fullUser.Location != nil {
				userProps = append(userProps, fmt.Sprintf("Location: %s", *fullUser.Location))
			}
			if fullUser.Email != nil {
				userProps = append(userProps, fmt.Sprintf("Email: %s", *fullUser.Email))
			}
			userData = append(userData, strings.Join(userProps, ", "))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	writeToFile(query, userData)
}

func fetchUsers(client *github.Client, userId int64) (*github.User, error) {
	user, _, err := client.Users.GetByID(context.Background(), userId)
	return user, err
}

func buildFileName() string {
	return time.Now().Format("20060102150405")
}

func MakeQueryToFilename(query string) (queryInFilename string) {
	replacer := strings.NewReplacer(":", "-", " ", "--", "#", "sharp")
	return replacer.Replace(query)
}

func writeToFile(query string, userData []string) {
	fileName := fmt.Sprintf("%s_%s_%s.txt", "result_", MakeQueryToFilename(query), buildFileName())
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range userData {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()
}
