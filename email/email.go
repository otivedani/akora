package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"time"

	"mime/multipart"
	"net/mail"
	"net/textproto"

	"gitlab.com/ameagaria.io/akora/config"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type encodedFile struct {
	fileName   string
	base64Data string
	mimeType   string
}

func ComposeEmail(ctx context.Context, targetEmail *mail.Address, attachmentFiles []string, isDraft bool) {

	// initialize gmail client

	client := config.GetGoogleClient(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// print content only, skip send/drafts
	// fmt.Println(createRawEmail(targetEmail, &attachmentFiles).String())

	var message gmail.Message
	message.Raw = base64.URLEncoding.EncodeToString(createRawEmail(targetEmail, &attachmentFiles).Bytes())

	if isDraft {
		var draft gmail.Draft
		draft.Message = &message

		_, err = srv.Users.Drafts.Create("me", &draft).Do()
		if err != nil {
			panic(err)
		}
		fmt.Println("Draft Created!")
	} else {
		_, err = srv.Users.Messages.Send("me", &message).Do()
		if err != nil {
			panic(err)
		}
		fmt.Println("Sent!")
	}

}

func createRawEmail(targetEmail *mail.Address, attachmentFiles *[]string) *bytes.Buffer {
	mixedContent := &bytes.Buffer{}
	mixedWriter := multipart.NewWriter(mixedContent)

	// create header
	header := make(textproto.MIMEHeader)
	header.Set("MIME-Version", "1.0")
	header.Set("Content-Type", "multipart/mixed; boundary="+mixedWriter.Boundary())
	header.Set("To", targetEmail.String())
	header.Set("From", config.GetSenderEmail())
	header.Set("Subject", "Halo")
	mixedWriter.CreatePart(header)

	// body message
	bodyHeader := make(textproto.MIMEHeader)
	bodyHeader.Set("Content-Type", "text/html; charset=UTF-8")
	bodyHeader.Set("Content-Transfer-Encoding", "7bit")
	mixedWriter.CreatePart(bodyHeader)
	mixedContent.WriteString("<h1>Hualo</h1>\n")

	// read files
	encodedFiles := make([]encodedFile, len(*attachmentFiles))
	for i, file := range *attachmentFiles {
		fb, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		var encFile = &encodedFiles[i]
		encFile.fileName = file
		encFile.mimeType = mime.TypeByExtension(filepath.Ext(file))
		encFile.base64Data = base64.StdEncoding.EncodeToString(fb)
	}

	for _, encFile := range encodedFiles {
		attchHeader := make(textproto.MIMEHeader)
		attchHeader.Set("Content-Type", fmt.Sprintf("%v; name='%v'",
			encFile.mimeType, encFile.fileName))
		attchHeader.Set("Content-Transfer-Encoding", "base64")
		attchHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename='%v'",
			encFile.fileName))
		mixedWriter.CreatePart(attchHeader)
		mixedContent.WriteString(encFile.base64Data)
	}
	// "Content-Type: " + fileMIMEType + "; name=" + string('"') + fileName + string('"') + " \n" +
	// "MIME-Version: 1.0\n" +
	// "Content-Transfer-Encoding: base64\n" +
	// "Content-Disposition: attachment; filename=" + string('"') + fileName + string('"') + " \n\n" +
	// chunkSplit(fileData, 76, "\n") +

	// fmt.Println(mixedContent.String())

	return mixedContent

}

// RetrieveSentMails dumps sent mail to db (and delete)
// TODO : separate delete operations
func RetrieveSentMails(ctx context.Context, f func(*SentMail) error) (err error) {

	client := config.GetGoogleClient(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	userId := "me"
	// userId := config.GetSenderEmail()
	// list := srv.Users.Messages.List("ameagaria.io@gmail.com")
	err = srv.Users.Messages.List(userId).Q(`
		from:(`+userId+`)
		is:sent
		(has:attachment OR has:drive)
		subject:(Permisi Paket!)
		`).
		Pages(ctx, func(lmr *gmail.ListMessagesResponse) error {
			for _, msg := range lmr.Messages {
				getMsg := srv.Users.Messages.Get(userId, msg.Id).Format("full")
				response, err := getMsg.Do()
				if err != nil || response == nil {
					log.Fatalf("error get call id:%+v", msg.Id)
				}

				mh := FromMessagePartHeader(response.Payload.Headers)

				for _, v := range response.Payload.Parts {
					if v.Filename != "" {
						mh.FileNames = append(mh.FileNames, v.Filename)
					}
				}

				mh.MessageId = msg.Id

				if err = f(mh); err != nil {
					return err
				}
			}
			return nil
		})

	if err != nil {
		return err
	}

	return nil
}

type SentMail struct {
	Date time.Time
	From mail.Address
	To   mail.Address

	FileNames []string

	MessageId string

	// MIMEVersion string       `name:"MIMEVersion"` //Value:1.0 ForceSendFields:[] NullFields:[]}
	// MessageID   string       `name:"MessageID"`
	// Subject     string       `name:"Subject"`
	// ContentType string       `name:"ContentType"` //Value:multipart/mixed; boundary="0000000000006bbc9d060acb9105" ForceSendFields:[] NullFields:[]}
}

func (mh *SentMail) ParseDate(s string) {
	t, _ := time.Parse(time.RFC1123Z, s)
	mh.Date = t
}

func (mh *SentMail) ParseFrom(s string) {
	addr, _ := mail.ParseAddress(s)
	if addr == nil {
		return
	}
	mh.From = *addr
}

func (mh *SentMail) ParseTo(s string) {
	addr, _ := mail.ParseAddress(s)
	if addr == nil {
		return
	}
	mh.To = *addr
}

// NewSentMail from gmail.MessagePartHeader
func FromMessagePartHeader(header []*gmail.MessagePartHeader) *SentMail {
	mh := &SentMail{}
	for _, h := range header {
		switch h.Name {
		case "Date":
			mh.ParseDate(h.Value)
		case "From":
			mh.ParseFrom(h.Value)
		case "To":
			mh.ParseTo(h.Value)
		}
	}
	return mh
}

// DeleteBatchMail
func DeleteBatchMail(ctx context.Context, msgids []string) {
	client := config.GetGoogleClient(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// max 1000 per request
	for i := 0; i < len(msgids)/1000; i++ {
		err = srv.Users.Messages.BatchDelete("me", &gmail.BatchDeleteMessagesRequest{
			Ids: msgids[i*1000 : (i+1)*1000],
		}).Do()

		if err != nil {
			log.Fatalf("delete error  %+v", err)
		}
	}

	fmt.Print(len(msgids), " deleted. bye")

}
