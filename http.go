package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// NewHttpCommand creates a new HTTP request command for the CLI.
func NewHttpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "http",
		Aliases: []string{"-http", "--http", "h", "-h", "--h"},
		Short: "\n------------------------------------\n" +
			"| COMMAND: http                    |\n" +
			"| TYPE: Http Request               |\n" +
			"| INPUT:                           |\n" +
			"|   1. [GET]  URL                  |\n" +
			"|   2. [POST] URL JSON             |\n" +
			"| EXAMPLES:                        |\n" +
			"|   1. 'https://a.com'             |\n" +
			"|   2. 'https://a.com' '{\"a\":1}'   |\n" +
			"------------------------------------\n",
		//Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var argLen = len(args)
			args0 := args[0]
			method := "GET"
			var body io.Reader
			var data string
			url := args0
			if argLen == 2 {
				method = "POST"
				data = args[1]

				// 接收的字符串双引号都没有了，需要处理
				data = strings.ReplaceAll(data, "{", `{"`)
				data = strings.ReplaceAll(data, "}", `"}`)
				data = strings.ReplaceAll(data, ":", `":"`)
				data = strings.ReplaceAll(data, ",", `","`)

				var jsonMap map[string]interface{}
				err := json.Unmarshal([]byte(data), &jsonMap)
				if err != nil {
					log.Fatal(err)
				}
				jsonData, _ := json.Marshal(jsonMap)
				body = bytes.NewBuffer(jsonData)
			}

			req, err := http.NewRequest(method, url, body)
			if err != nil {
				log.Fatal(err)
			}
			if method == "POST" {
				req.Header.Set("Content-Type", "application/json")
			}

			startTime := time.Now()
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			endTime := time.Now()
			duration := endTime.Sub(startTime).Milliseconds()

			fmt.Println("")
			fmt.Printf("Request Url       ：%s\n", url)
			if method == "POST" {
				fmt.Printf("Request Post json : %s\n", data)
			}

			index := 0
			for name, values := range resp.Header {
				for _, value := range values {
					if index == 0 {
						fmt.Printf("Response Headers  : %s: %s\n", name, value)
					} else {
						fmt.Printf("                    %s: %s\n", name, value)
					}
				}
				index++
			}
			fmt.Printf("Response Start    : %s\n", startTime.Format("2006-01-02 15:04:05.000"))
			fmt.Printf("Response End      : %s\n", endTime.Format("2006-01-02 15:04:05.000"))
			fmt.Printf("Response Duration : %d ms\n", duration)
			fmt.Printf("Response Status   : %d\n", resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Response Body     : %s\n", string(respBody))
			fmt.Println("")
		},
	}
	return cmd
}
