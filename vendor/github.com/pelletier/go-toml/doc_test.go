// code examples for godoc

package toml_test

import (
	"fmt"
	"log"

	toml "github.com/pelletier/go-toml"
)

func Example_tree() ***REMOVED***
	config, err := toml.LoadFile("config.toml")

	if err != nil ***REMOVED***
		fmt.Println("Error ", err.Error())
	***REMOVED*** else ***REMOVED***
		// retrieve data directly
		user := config.Get("postgres.user").(string)
		password := config.Get("postgres.password").(string)

		// or using an intermediate object
		configTree := config.Get("postgres").(*toml.Tree)
		user = configTree.Get("user").(string)
		password = configTree.Get("password").(string)
		fmt.Println("User is", user, " and password is", password)

		// show where elements are in the file
		fmt.Printf("User position: %v\n", configTree.GetPosition("user"))
		fmt.Printf("Password position: %v\n", configTree.GetPosition("password"))
	***REMOVED***
***REMOVED***

func Example_unmarshal() ***REMOVED***
	type Employer struct ***REMOVED***
		Name  string
		Phone string
	***REMOVED***
	type Person struct ***REMOVED***
		Name     string
		Age      int64
		Employer Employer
	***REMOVED***

	document := []byte(`
	name = "John"
	age = 30
	[employer]
		name = "Company Inc."
		phone = "+1 234 567 89012"
	`)

	person := Person***REMOVED******REMOVED***
	toml.Unmarshal(document, &person)
	fmt.Println(person.Name, "is", person.Age, "and works at", person.Employer.Name)
	// Output:
	// John is 30 and works at Company Inc.
***REMOVED***

func ExampleMarshal() ***REMOVED***
	type Postgres struct ***REMOVED***
		User     string `toml:"user"`
		Password string `toml:"password"`
		Database string `toml:"db" commented:"true" comment:"not used anymore"`
	***REMOVED***
	type Config struct ***REMOVED***
		Postgres Postgres `toml:"postgres" comment:"Postgres configuration"`
	***REMOVED***

	config := Config***REMOVED***Postgres***REMOVED***User: "pelletier", Password: "mypassword", Database: "old_database"***REMOVED******REMOVED***
	b, err := toml.Marshal(config)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	fmt.Println(string(b))
	// Output:
	// # Postgres configuration
	// [postgres]
	//
	//   # not used anymore
	//   # db = "old_database"
	//   password = "mypassword"
	//   user = "pelletier"
***REMOVED***

func ExampleUnmarshal() ***REMOVED***
	type Postgres struct ***REMOVED***
		User     string
		Password string
	***REMOVED***
	type Config struct ***REMOVED***
		Postgres Postgres
	***REMOVED***

	doc := []byte(`
	[postgres]
	user = "pelletier"
	password = "mypassword"`)

	config := Config***REMOVED******REMOVED***
	toml.Unmarshal(doc, &config)
	fmt.Println("user=", config.Postgres.User)
	// Output:
	// user= pelletier
***REMOVED***
