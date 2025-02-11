package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/UltimateTournament/go-structfilter/structfilter"
)

// User represents a user entry in a user database.
type User struct {
	Name          string
	Password      string `req_role:"admin superadmin"`
	PasswordAdmin string `req_role:"superadmin"`
	LoginTime     time.Time
	*EmbeddedTest `req_role:"admin superadmin"`
	AgeP          *Duration
	Age           Duration
}

type EmbeddedTest struct {
	EmbeddedField string
}

var age = Duration(123 * time.Hour)

var userDB = map[string][]*User{
	"foo_DB": {
		{
			Name:          "Alice",
			Password:      "$6$sensitive",
			PasswordAdmin: "$6$verysensitive",
			LoginTime:     time.Now().Add(-24 * time.Hour),
			Age:           Duration(999 * time.Hour),
		},
		{
			Name:      "Bob",
			Password:  "$6$private",
			LoginTime: time.Now().Add(-36 * time.Hour),
			EmbeddedTest: &EmbeddedTest{
				EmbeddedField: "test1",
			},
			AgeP: &age,
		},
	},
}

func main() {
	userRole := "admin"
	converted, err := createRoleStructFilter(userRole).Convert(userDB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", converted)
	jsonData, err := json.MarshalIndent(converted, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

var jsonMarshallerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()

func createRoleStructFilter(userRole string) *structfilter.T {
	roleFilter := func(f *structfilter.Field) error {
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
	}
	keepCustomJsonFilter := func(f *structfilter.Field) error {
		hasCustomJsonMarshaller := f.Type.AssignableTo(jsonMarshallerType)
		if hasCustomJsonMarshaller {
			f.KeepRaw()
		}
		return nil
	}
	filter := structfilter.New(roleFilter, keepCustomJsonFilter)
	return filter
}
