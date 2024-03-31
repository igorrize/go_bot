package ud_client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type TLDRClient interface {
	ShortDefinitions(longData string) error
}

type ShortClient struct {
	apiKey string
	host   string
}

func NewShortClient(apiKey, host string) *ShortClient {
	return &ShortClient{
		apiKey: apiKey,
		host: host,
	}
}

func (c *ShortClient) ShortDefinitions(longData string) (string, error) {

	url := "https://tldrthis.p.rapidapi.com/v1/model/abstractive/summarize-text/"

	payload := strings.NewReader(fmt.Sprintf(`{"text": "%s"}`, longData))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", c.apiKey)
	req.Header.Add("X-RapidAPI-Host", "tldrthis.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, _ := io.ReadAll(res.Body)
	var data map[string]interface{}

	err := json.Unmarshal(body, &data)
	if err != nil {
		return "err", nil
	}

	summary, ok := data["summary"].(string)
	if !ok {
		return "err", nil
	}
	return summary, nil

}

