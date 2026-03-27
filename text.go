package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

// NewTextCommand creates a new text encoding/decoding command for the CLI.
func NewTextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "text",
		Aliases: []string{"t", "-t", "--t", "-text", "--text"},
		Short: "\n-------------------------------------\n" +
			"| COMMAND: text                   |\n" +
			"| TYPE: Text Encode/Decode       |\n" +
			"| INPUT:                          |\n" +
			"|   1. <text>: convert all       |\n" +
			"|   2. base64 encode <text>      |\n" +
			"|   3. base64 decode <text>      |\n" +
			"|   4. url encode <text>         |\n" +
			"|   5. url decode <text>         |\n" +
			"|   6. md5 <text>                |\n" +
			"| EXAMPLES:                       |\n" +
			"|   1. hello world               |\n" +
			"|   2. base64 encode hello       |\n" +
			"|   3. base64 decode aGVsbG8=    |\n" +
			"|   4. url encode hello world    |\n" +
			"|   5. url decode hello%20world  |\n" +
			"|   6. md5 password              |\n" +
			"-------------------------------------\n",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("invalid arguments, use 'text --help' for usage")
			}

			operation := strings.ToLower(args[0])

			if len(args) >= 2 {
				subOperation := strings.ToLower(args[1])
				var text string

				if len(args) > 2 {
					text = strings.Join(args[2:], " ")
				}

				switch operation {
				case "base64":
					return handleBase64(subOperation, text)
				case "url":
					return handleURL(subOperation, text)
				case "md5":
					return handleMD5(args[1:])
				default:
					text = strings.Join(args, " ")
					return handleAllConversions(text)
				}
			}

			return handleAllConversions(operation)
		},
	}
	return cmd
}

// handleAllConversions performs all conversions on the given text.
func handleAllConversions(text string) error {
	fmt.Printf("Input: %s\n\n", text)

	base64Encoded := base64.StdEncoding.EncodeToString([]byte(text))
	fmt.Printf("✓ Base64 Encode: %s\n", base64Encoded)

	base64Decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		fmt.Printf("✗ Base64 Decode: failed (%v)\n", err)
	} else {
		fmt.Printf("✓ Base64 Decode: %s\n", string(base64Decoded))
	}

	urlEncoded := url.QueryEscape(text)
	fmt.Printf("✓ URL Encode: %s\n", urlEncoded)

	urlDecoded, err := url.QueryUnescape(text)
	if err != nil {
		fmt.Printf("✗ URL Decode: failed (%v)\n", err)
	} else {
		fmt.Printf("✓ URL Decode: %s\n", urlDecoded)
	}

	hash := md5.Sum([]byte(text))
	fmt.Printf("✓ MD5: %x\n", hash)

	return nil
}

// handleBase64 handles base64 encoding and decoding operations.
func handleBase64(operation, text string) error {
	switch operation {
	case "encode":
		if text == "" {
			return fmt.Errorf("text is required for base64 encode")
		}
		encoded := base64.StdEncoding.EncodeToString([]byte(text))
		fmt.Printf("input  : %s\n", text)
		fmt.Printf("encode : %s\n", encoded)
		return nil

	case "decode":
		if text == "" {
			return fmt.Errorf("text is required for base64 decode")
		}
		decoded, err := base64.StdEncoding.DecodeString(text)
		if err != nil {
			log.Fatal("base64 decode error:", err)
		}
		fmt.Printf("input  : %s\n", text)
		fmt.Printf("decode : %s\n", string(decoded))
		return nil

	default:
		return fmt.Errorf("unknown base64 operation: %s", operation)
	}
}

// handleURL handles URL encoding and decoding operations.
func handleURL(operation, text string) error {
	switch operation {
	case "encode":
		if text == "" {
			return fmt.Errorf("text is required for url encode")
		}
		encoded := url.QueryEscape(text)
		fmt.Printf("input  : %s\n", text)
		fmt.Printf("encode : %s\n", encoded)
		return nil

	case "decode":
		if text == "" {
			return fmt.Errorf("text is required for url decode")
		}
		decoded, err := url.QueryUnescape(text)
		if err != nil {
			log.Fatal("url decode error:", err)
		}
		fmt.Printf("input  : %s\n", text)
		fmt.Printf("decode : %s\n", decoded)
		return nil

	default:
		return fmt.Errorf("unknown url operation: %s", operation)
	}
}

// handleMD5 handles MD5 hashing operations.
func handleMD5(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("text is required for md5")
	}
	text := strings.Join(args, " ")
	hash := md5.Sum([]byte(text))
	fmt.Printf("input : %s\n", text)
	fmt.Printf("md5   : %x\n", hash)
	return nil
}
