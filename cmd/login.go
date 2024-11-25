/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/cingram16/spotify/internal/spotify"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login as spotify user",
	Long:  "login as spotify user",
	Run: func(cmd *cobra.Command, args []string) {
		spotify.NewClient()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
