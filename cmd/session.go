/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"time"
	"ttm/pkg/fs"
	"ttm/pkg/models"
	"ttm/pkg/render"

	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage sessions",
	Args:  cobra.MinimumNArgs(1),
}

var startCmd = &cobra.Command{
	Use:   "start [task_id]",
	Short: "Start a new session",
	Args:  cobra.MinimumNArgs(1),
	Run:   startHandler,
}

var endCmd = &cobra.Command{
	Use:   "end",
	Short: "End a session",
	Run:   endHandler,
}

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a session",
	Run:   cancelHandler,
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get session info",
	Run:   infoHandler,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sessions",
	Run:   listHandler,
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	sessionCmd.AddCommand(startCmd)
	sessionCmd.AddCommand(endCmd)
	sessionCmd.AddCommand(cancelCmd)
	sessionCmd.AddCommand(infoCmd)
	sessionCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("id", "i", "", "Filter by task ID")
	listCmd.Flags().StringP("days", "d", "", "Previous days to show")
}

func startHandler(cmd *cobra.Command, args []string) {
	if fs.SessionFileExists() {
		fmt.Println("Session already started. Please end the current session first.")
		return
	}

	taskId := args[0]

	_, err := fs.CreateSessionFile(taskId, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}

	render.RenderSessionStart(taskId)
}

func endHandler(cmd *cobra.Command, args []string) {
	if !fs.SessionFileExists() {
		fmt.Println("No session found. Please start a session first.")
		return
	}

	sf, err := fs.RemoveSessionFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	taskId, err := strconv.Atoi(sf.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	taskStore.AddSession(models.Session{
		TaskId:    int64(taskId),
		StartTime: sf.StartTime,
		EndTime:   time.Now(),
	})

	render.RenderSessionEnd(sf)
}

func cancelHandler(cmd *cobra.Command, args []string) {
	if !fs.SessionFileExists() {
		fmt.Println("No session found. Please start a session first.")
		return
	}

	_, err := fs.RemoveSessionFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	render.RenderSessionCancel()
}

func infoHandler(cmd *cobra.Command, args []string) {
	if !fs.SessionFileExists() {
		fmt.Println("No session found. Please start a session first.")
		return
	}

	sf, err := fs.ReadSessionFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	render.RenderSessionInfo(sf)
}
