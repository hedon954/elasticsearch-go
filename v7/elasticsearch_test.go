package elasticsearch

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

func TestClient_Version(t *testing.T) {
	client := DefaultClient()
	t.Log(client.Version())
}
