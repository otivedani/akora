package ingest

import (
	"context"
	"testing"
)

func TestIngestData(t *testing.T) {
	ctx := context.Background()
	IngestData(ctx)
}
