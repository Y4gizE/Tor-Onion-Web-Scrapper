package cmd

import (
	"fmt"
	"webscraper/internal"

	"github.com/spf13/cobra"
)

var (
	url            string
	fetchHtml      bool
	fetchLinks     bool
	takeScreenshot bool
)

var rootCmd = &cobra.Command{
	Use:   "webscraper",
	Short: "A simple web scraper for Dark Web sites",
	Long:  `This tool scrapes content from Dark Web sites for research purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		torSession, err := internal.StartTorSession() // Start the Tor session
		if err != nil {
			fmt.Println("Failed to create Tor session:", err)
			return
		}
		defer torSession.Stop() // Clean up Tor session after use

		if fetchHtml {
			internal.FetchAndSaveHTML(url, torSession) // Pass the torSession
		}
		if fetchLinks {
			internal.FetchAndSaveLinks(url, torSession) // Pass the torSession
		}
		if takeScreenshot {
			internal.TakeScreenshotViaTor(url, torSession) // Pass the torSession
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the Dark Web site (must end with .onion)")
	rootCmd.Flags().BoolVar(&fetchHtml, "html", false, "Fetch HTML content from the given URL")
	rootCmd.Flags().BoolVar(&fetchLinks, "links", false, "Fetch links from the given URL")
	rootCmd.Flags().BoolVar(&takeScreenshot, "screenshot", false, "Take a screenshot of the given URL")
	rootCmd.MarkFlagRequired("url")
}
