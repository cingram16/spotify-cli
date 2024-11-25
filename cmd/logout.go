/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/cingram16/spotify/internal/spotify"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout spotify user",
	Long:  "Logout spotify user",
	Run: func(cmd *cobra.Command, args []string) {
		err := spotify.RemoveConfig()
		if err != nil {
			fmt.Println("Error logging out")
		}
		fmt.Println("logout called")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

}
