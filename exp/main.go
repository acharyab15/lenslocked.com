package main

import (
	"html/template"
	"os"
)

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	data := map[string]string{
		"Name": "John Smith",
		"Age":  "42",
	}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
