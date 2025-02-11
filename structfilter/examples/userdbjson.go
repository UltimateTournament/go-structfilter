package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/UltimateTournament/go-structfilter/structfilter"
)

// User represents a user entry in a user database.
type User struct {
	Name          string
	Password      string `req_role:"admin superadmin"`
	PasswordAdmin string `req_role:"superadmin"`
	LoginTime     int64
}

var userDB = map[string][]*User{
	"foo_DB": {
		{
			Name:          "Alice",
			Password:      "$6$sensitive",
			PasswordAdmin: "$6$verysensitive",
			LoginTime:     1234567890,
		},
		{
			Name:      "Bob",
			Password:  "$6$private",
			LoginTime: 1357924680,
		},
	},
}

func main() {
	userRole := "editor"
	converted, err := createRoleStructFilter(userRole).Convert(userDB)
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(converted, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

func createRoleStructFilter(userRole string) *structfilter.T {
	filter := structfilter.New(
		func(f *structfilter.Field) error {
			reqRolesStr := f.Tag.Get("req_role")
			if reqRolesStr == "" {
				return nil
			}
			reqRoles := strings.Split(reqRolesStr, " ")
			for _, reqRole := range reqRoles {
				if userRole == reqRole {
					return nil
				}
			}
			f.Remove()
			return nil
		},
	)
	return filter
}
