/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	useBeta bool
)

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update y360c app",
	Long:  `Use this command to update y360c to the latest version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("checking the latest version...")

		wd, _ := os.Getwd()
		filename := fmt.Sprintf("y360c-%s-%s.zip", strings.ToLower(runtime.GOOS), strings.ToLower(runtime.GOARCH))
		url := fmt.Sprintf("https://github.com/foxsoft2005/y360c/releases/latest/download/%s", filename)

		req, _ := http.NewRequest("GET", url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln("failed to get update:", err)
		}
		defer resp.Body.Close()

		f, _ := os.OpenFile(filepath.Join(wd, filename), os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()

		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"downloading",
		)
		io.Copy(io.MultiWriter(f, bar), resp.Body)
		fmt.Printf("update downloaded to %v, extract it manually over the existing files", filepath.Join(wd, filename))
	},
}

func init() {
	UpdateCmd.Flags().BoolVarP(&useBeta, "beta", "b", false, "use beta version for update")
}
