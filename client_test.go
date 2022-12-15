package elasticsearch_go

import (
	"testing"
)

func TestClient_Ping(t *testing.T) {
	client := DefaultClient()
	if client.Ping() {
		t.Log("ping es successfullt")
	} else {
		t.Error("ping es failed")
	}
}
