/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/cingram16/spotify/internal/spotify"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	spfy "github.com/zmb3/spotify/v2"
	"log"
	"strings"
)

// playlistCmd represents the playlist command
var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "Interact with spotify playlists",
	Long: `Interact with spotify playlists. For example: 

			List all playlists:
			spotify-cli playlist ls`,
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all playlists",
	Long:  "List all playlists of the authenticated Spotify user",
	Run: func(cmd *cobra.Command, args []string) {
		client := spotify.NewClient()

		playlists, err := client.ListPlaylists(cmd.Context())
		if err != nil {
			log.Fatalf("Error listing playlists: %v", err)
		}

		for _, playlist := range playlists {
			fmt.Println(playlist.Name)
		}
	},
}

// rotateCmd represents the rotate command
var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate a playlist",
	Long:  "Rotate a playlist with new tracks based on user's top artists and tracks",
	Run: func(cmd *cobra.Command, args []string) {
		client := spotify.NewClient()

		playlists, err := client.ListPlaylists(cmd.Context())
		if err != nil {
			log.Fatalf("Error listing playlists: %v", err)
		}

		playlistNames := make([]string, len(playlists))
		for i, playlist := range playlists {
			playlistNames[i] = playlist.Name
		}

		prompt := promptui.Select{
			Label: "Select Playlist",
			Items: playlistNames,
		}

		_, result, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		var selectedPlaylist *spfy.SimplePlaylist
		for _, playlist := range playlists {
			if playlist.Name == result {
				selectedPlaylist = playlist
				break
			}
		}

		if selectedPlaylist == nil {
			log.Fatalf("Selected playlist not found")
		}

		promptArtist := promptui.Prompt{
			Label: "Enter comma-separated artist names to seed",
		}
		artistNames, err := promptArtist.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		artistNameList := strings.Split(artistNames, ",")
		var artistIDs []spfy.ID
		for _, artistName := range artistNameList {
			fmt.Println("Searching for artist", artistName)
			artist, err := client.GetFirstSearchArtist(cmd.Context(), strings.TrimSpace(artistName))
			fmt.Println("Artist found:", artist.Name)
			if err != nil {
				log.Fatalf("Error searching for artist %s: %v", artistName, err)
			}
			if artist == nil {
				log.Fatalf("Artist %s not found", artistName)
			}
			artistIDs = append(artistIDs, artist.ID)
		}

		err = client.RotatePlaylist(
			cmd.Context(),
			*selectedPlaylist,
			spfy.Seeds{
				Artists: artistIDs,
			})
		if err != nil {
			log.Fatalf("Error rotating playlist: %v", err)
		}

		fmt.Println("Playlist rotated successfully")
	},
}

func init() {
	rootCmd.AddCommand(playlistCmd)
	playlistCmd.AddCommand(listCmd)
	playlistCmd.AddCommand(rotateCmd)
}
