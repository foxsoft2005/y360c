/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package internal

import (
	"log"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/spf13/cobra"
)

var (
	useBeta bool
)

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates y360c app",
	Long:  `Use this command to update y360c application to the actual version.`,
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://mlgmow01app010.linru.grp/files/y360c.zip"

		client := grab.NewClient()
		req, _ := grab.NewRequest(".", url)

		if useBeta {
			log.Println("Downloading the beta version...")
		} else {
			log.Println("Downloading the actual version...")
		}

		resp := client.Do(req)
		log.Printf("  %v", resp.HTTPResponse.Status)

		t := time.NewTicker(200 * time.Millisecond)
		defer t.Stop()

	Loop:
		for {
			select {
			case <-t.C:
				log.Printf("  transferred %v / %v bytes (%.2f%%)", resp.BytesComplete(), resp.Size(), 100*resp.Progress())
			case <-resp.Done:
				break Loop
			}
		}

		if err := resp.Err(); err != nil {
			log.Fatalln("Download failed:", err)
		}

		log.Printf("Update saved to ./%v, please extract it over the existing version manually", resp.Filename)
	},
}

func init() {
	UpdateCmd.Flags().BoolVarP(&useBeta, "beta", "b", false, "Use beta version for update")
}
