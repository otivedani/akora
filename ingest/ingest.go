package ingest

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"gitlab.com/ameagaria.io/akora/config"
	"gitlab.com/ameagaria.io/akora/model"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var spreadsheetId string = os.Getenv("SPREADSHEET_ID")

func IngestData(ctx context.Context) {

	client := config.GetGoogleClient(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// read from db
	startRow := 2
	readRange := fmt.Sprintf("Form Responses 1!A%d:D", startRow)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			if len(row) > 0 {
				d := newPostcard(row)
				// store to db
				fmt.Printf("%+v\n", d)
			}
		}
		log.Printf("INFO: %d entry was added.", len(resp.Values))
	}
}

// Populate get cell results into new model.Postcard
func newPostcard(row []interface{}) *model.Postcard {
	postcard := model.Postcard{}
	execorder := []func(*model.Postcard, interface{}){
		parseTimestamp, // cell A
		parseOrderNo,   // cell B
		parseEmail,     // cell C
		parsePhone,     // cell D
	}
	for i, v := range row {
		execorder[i](&postcard, v)
	}

	return &postcard
}

func parseOrderNo(postcard *model.Postcard, value interface{}) {
	postcard.SetOrderNo(fmt.Sprint(value))
}

func parseEmail(postcard *model.Postcard, value interface{}) {
	postcard.SetEmail(fmt.Sprint(value))
}

func parseTimestamp(postcard *model.Postcard, value interface{}) {
	res, err := time.Parse("1/2/2006 15:04:05", fmt.Sprint(value))
	if err != nil {
		log.Fatalf("Unable to parse data from sheet: %v", err)
	}
	postcard.SetTimestamp(res.Unix())
}

func parsePhone(postcard *model.Postcard, value interface{}) {
	re := regexp.MustCompile(`^(?:0|\+?62)|([\s-])`).ReplaceAllString(fmt.Sprint(value), "")
	if re == "" {
		postcard.SetPhone(0)
		return
	}
	num, err := strconv.ParseUint("62"+re, 10, 64)
	if err != nil {
		log.Fatalf("Unable to parse data from sheet: %v", err)
	}
	postcard.SetPhone(uint64(num))
}
