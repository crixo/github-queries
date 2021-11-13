// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The simple command demonstrates a simple functionality which
// prompts the user for a GitHub username and lists all the public
// organization memberships of the specified username.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v40/github"

	"golang.org/x/oauth2"
)

// Fetch all the public organizations' membership of a user.
//
func fetchOrganizations(username string) ([]*github.Repository, error) {
	client := github.NewClient(nil)
	orgs, _, err := client.Repositories.List(context.Background(), username, nil)
	return orgs, err
}

// func fetchUsers(username string) ([]*github.User, error) {
// 	client := github.NewClient(nil)
// 	opts :=
// 	users, _, err := client.Users.ListAll(context.Background(), opts)
// 	return users, err
// }

func main() {
	//fmt.Print("Enter User Location : ")
	//var input string
	// fmt.Scanln(&input)
	// fmt.Print(input)
	//input = "Italy"

	ctx := context.Background()

	var token string = os.Getenv("PAT")

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
	for {

		// if lastId > 0 {
		// 	opt.Since = int64(lastId)
		// }

		// https://docs.github.com/en/rest/reference/search#search-users
		query := "language:go location:italy"

		userSearch, resp, err := client.Search.Users(ctx, query, opts)
		if err != nil {
			//return err
		}

		for _, user := range userSearch.Users {
			allUsers = append(allUsers, user)
			fmt.Println(*user.Login)

			//lastId = int(*user.ID)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
}
