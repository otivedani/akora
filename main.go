package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"

	"gitlab.com/ameagaria.io/akora/email"
)

func main() {
	// test, intend to run `go ...`
	ctx := context.Background()
	// ingest.IngestData(ctx)

	// type StudioSentMail struct {
	// 	Date      time.Time    `db:"date"`
	// 	From      mail.Address `db:"from"`
	// 	To        mail.Address `db:"to"`
	// 	FileNames []string     `db:"filenames"`
	// 	MessageId string       `db:"gmail__message_id"`
	// }

	// db := config.GetSqliteClient()

	// _, err := db.Exec(`
	// CREATE TABLE IF NOT EXISTS studio_sent_mail (
	// 	"date" DATETIME,
	// 	"from" TEXT,
	// 	"to" TEXT,
	// 	"filenames" TEXT,
	// 	"gmail__message_id" TEXT
	// 	);
	// `)

	// if err != nil {
	// 	log.Fatalf("insert error create %+v", err)
	// }

	// email.RetrieveSentMails(ctx, func(mh *email.SentMail) error {
	// 	// data := StudioSentMail(*mh)
	// 	fns, _ := json.Marshal(mh.FileNames)

	// 	res, err := db.Exec(`
	// 	INSERT INTO studio_sent_mail (
	// 		"date","from","to","filenames","gmail__message_id"
	// 	) VALUES (?,?,?,?,?);
	// 	`,
	// 		mh.Date, mh.From.String(), mh.To.String(), string(fns), mh.MessageId,
	// 	)

	// 	if err != nil {
	// 		log.Fatalf("insert error exec %+v", err)
	// 	}

	// 	var id int64
	// 	if id, err = res.LastInsertId(); err != nil {
	// 		log.Fatalf("insert error row %+v", err)
	// 	}
	// 	fmt.Printf("%+v %v\n", mh, id)
	// 	return nil
	// })

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
