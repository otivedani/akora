package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"

	"gitlab.com/ameagaria.io/akora/email"
	"gitlab.com/ameagaria.io/akora/ingest"
)

func main() {
	// test, intend to run `go ...`
	ctx := context.Background()
	ingest.IngestData(ctx)

	// args and flags handling
	var isDraft = true
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "draft":
			isDraft = true
		case "send":
			isDraft = false
		default:
			log.Fatal("Either `draft` or `send` must be supplied.")
		}
	} else {
		log.Fatal("Either `draft` or `send` must be supplied.")
	}

	var rawEmail string
	if len(os.Args) > 2 {
		rawEmail = os.Args[2]
	} else {
		fmt.Printf("[0] Destination email :\n")
		if _, err := fmt.Scanln(&rawEmail); err != nil {
			log.Fatalf("Error on email input: %v", err)
		}
	}

	toEmail, err := mail.ParseAddress(rawEmail)
	if err != nil {
		log.Fatal(err)
	}

	var attachmentFiles []string
	var tempFile string
	if len(os.Args) > 3 {
		attachmentFiles = os.Args[3:]
	} else {
		fmt.Printf("[1] Attachment files :\n")
		for {
			_, err := fmt.Scanln(&tempFile)
			if err != nil {
				log.Fatalf("Error on files input: %v", err)
			}
			attachmentFiles = append(attachmentFiles, tempFile)
		}
	}

	email.ComposeEmail(ctx, toEmail, attachmentFiles, isDraft)
}
