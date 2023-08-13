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

	return
}

func createRawEmail(targetEmail *mail.Address, attachmentFiles *[]string) *bytes.Buffer {
	mixedContent := &bytes.Buffer{}
	mixedWriter := multipart.NewWriter(mixedContent)

	// create header
	header := make(textproto.MIMEHeader)
	header.Set("MIME-Version", "1.0")
	header.Set("Content-Type", "multipart/mixed; boundary="+mixedWriter.Boundary())
	header.Set("To", targetEmail.String())
	header.Set("From", os.Getenv("SENDER_EMAIL"))
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
