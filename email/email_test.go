package email

import (
	"context"
	"fmt"
	"testing"
)

func TestRetrieveSentMails(t *testing.T) {
	ctx := context.Background()

	RetrieveSentMails(ctx, func(mh *SentMail) error {
		fmt.Printf("%+v\n", mh)
		return nil
	})
}
