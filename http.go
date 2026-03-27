package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func isHTTPURL(input string) bool {
	return strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")
}

func handleHTTPRequest(input string) error {
	input = strings.TrimSpace(input)
	input = strings.Trim(input, `"'`)

	var url string
	var data string

	var urlEnd = strings.Index(input, " ")
	if urlEnd == -1 {
		url = input
	} else {
		url = input[:urlEnd]
		data = strings.TrimSpace(input[urlEnd+1:])
		data = strings.Trim(data, `"'`)
	}

	method := "GET"
	var body io.Reader

	if data != "" {
		method = "POST"
		var jsonMap map[string]interface{}
		err := json.Unmarshal([]byte(data), &jsonMap)
		if err != nil {
			return fmt.Errorf("invalid JSON data: %v", err)
		}
		jsonData, _ := json.Marshal(jsonMap)
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	startTime := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	endTime := time.Now()
	duration := endTime.Sub(startTime).Milliseconds()

	fmt.Println("")
	fmt.Printf("Request Url       ：%s\n", url)
	if method == "POST" {
		fmt.Printf("Request Post json : %s\n", data)
	}
	fmt.Printf("Request Method    ：%s\n", method)

	header := resp.Header
	for key, value := range header {
		fmt.Printf("Response Header %s: %s\n", key, strings.Join(value, ","))
	}
	fmt.Printf("Response status   ：%s\n", resp.Status)
	fmt.Printf("Response Time(ms) ：%dms\n", duration)

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Response Body:\n%s\n", string(bodyData))
	return nil
}

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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("URL is required")
			}
			input := strings.Join(args, " ")
			return handleHTTPRequest(input)
		},
	}
	return cmd
}

