package repository

import (
	"fmt"
	"log"

	"gitlab.com/ameagaria.io/akora/config"
)

func GetMailIdsBefore(days int) (messageIds []string) {

	sqliteDB := config.GetSqliteDB()

	daysMod := fmt.Sprintf(`-%d day`, days)
	err := sqliteDB.Select(&messageIds,
		`
			SELECT gmail__message_id FROM studio_sent_gmail 
			WHERE date < DATETIME('now', ?)
		`,
		daysMod,
	)

	if err != nil {
		log.Fatalf("repository.akora.GetMailIdsBefore:select::%+v", err)
	}

	return messageIds
}
