# go_telesign
telesign golang SDK

## Installation

`go get -u github.com/yedamao/go_telesign`

## Quick Start

```go
custom_id := "id"
api_key := "key"
client := New("https://rest-ww.telesign.com", custom_id, api_key)

resp, err := client.Send("xxxxxxxxxxxxx", "3388", "en-US", "hello, your code is $$CODE$$ .")
if err != nil {
    log.Fatal(err)
}
log.Println(resp)
```
