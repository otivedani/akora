package repository

import (
	"testing"
)

func TestGetMailIdsBefore(t *testing.T) {

	msgIds := GetMailIdsBefore(30)
	t.Logf(`%+v`, msgIds)
}
