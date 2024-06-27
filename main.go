/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ttm/cmd"
	_ "ttm/cmd/session" // Importing the session package so go build finds it
	_ "ttm/cmd/task"    // Importing the task package so go build finds it
)

func main() {
	cmd.Execute()
}
