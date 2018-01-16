/*
implement Telesign golang SDK

author: huanhuan8

*/

package telesign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

type RestClient struct {
	customer_id   string // Your customer_id string associated with your account.
	api_key       string //  Your api_key string associated with your account.
	rest_endpoint string //
	baseURL       *url.URL
}

type Status struct {
	UpdatedOn   string `json:"updated_on"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Reponse struct {
	ReferenceId string  `json:"reference_id"`
	ResourceUrl string  `json:"resource_uri"`
	SubResponse string  `json:"sub_resource"`
	Status      Status  `json:"status"`
	Errors      []Error `json:"errors"`
}

func New(customer_id, api_key string) *RestClient {
	return &RestClient{
		customer_id:   customer_id,
		api_key:       api_key,
		rest_endpoint: "https://rest-ww.telesign.com",
	}
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))

}

func (c *RestClient) generateTelesignHeaders(
	method,
	resource,
	url_encoded_fields,
	date_rfc2616,
	nonce,
	user_agent string) http.Header {

	h := http.Header{}

	if "" == date_rfc2616 {
		location, _ := time.LoadLocation("GMT")
		date_rfc2616 = time.Now().In(location).Format(time.RFC1123)
	}

	if "" == nonce {
		nonce = uuid.Must(uuid.NewV4()).String()
	}

	content_type := ""
	if "POST" == method || "PUT" == method {
		content_type = "application/x-www-form-urlencoded"
	}

	auth_method := "HMAC-SHA256"

	string_to_sign_builder := method + "\n" +
		content_type + "\n" +
		date_rfc2616 + "\n" +
		"x-ts-auth-method:" + auth_method + "\n" +
		"x-ts-nonce:" + nonce

	if content_type != "" && url_encoded_fields != "" {
		string_to_sign_builder += "\n" + url_encoded_fields
	}

	string_to_sign_builder += "\n" + resource

	api_key, _ := base64.StdEncoding.DecodeString(c.api_key)
	signature := computeHmac256(string_to_sign_builder, string(api_key))

	authorization := "TSA " + c.customer_id + ":" + signature

	// build header
	h.Set("Authorization", authorization)
	h.Set("Date", date_rfc2616)
	h.Set("Content-Type", content_type)
	h.Set("x-ts-auth-method", auth_method)
	h.Set("x-ts-nonce", nonce)

	if user_agent != "" {
		h.Set("User-Agent", user_agent)
	}

	return h
}

func (c *RestClient) newRequest(method, resource string, params url.Values) (*http.Request, error) {
	var req *http.Request
	var err error

	if "POST" == method || "PUT" == method {
		req, err = http.NewRequest(method, c.rest_endpoint+resource, strings.NewReader(params.Encode()))
	} else {
		req, err = http.NewRequest(method, c.rest_endpoint+resource, nil)
	}
	if err != nil {
		return nil, err
	}

	header := c.generateTelesignHeaders(method, resource, params.Encode(), "", "", "")

	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	return req, nil
}

func (c *RestClient) post(resource string, params url.Values) (*Reponse, error) {

	req, err := c.newRequest("POST", resource, params)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)

	var rep = new(Reponse)
	err = json.Unmarshal(data, &rep)
	if err != nil {
		return nil, err
	}

	return rep, nil
}

func (c *RestClient) get(resource string, params url.Values) error {
	// todo
	return nil
}

func (c *RestClient) put(resource string, params url.Values) error {
	// todo
	return nil
}

func (c *RestClient) delete(resource string, params url.Values) error {
	// todo
	return nil
}

func (c *RestClient) Send(to, code, lang, template string) (*Reponse, error) {

	content := url.Values{
		"phone_number": {to},
		"verify_code":  {code},
		"language":     {lang},
		"template":     {template},
	}

	return c.post("/v1/verify/sms", content)
}
