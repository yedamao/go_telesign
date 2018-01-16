package telesign

import (
	"log"
	"testing"
)

func newClient() *RestClient {
	custom_id := "id"
	api_key := "key"
	return New(custom_id, api_key)
}

func TestSend(t *testing.T) {
	client := newClient()

	resp, err := client.Send("xxxxxxxxxxxxx", "3388", "en-US", "haha your code is $$CODE$$ .")
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)
}
