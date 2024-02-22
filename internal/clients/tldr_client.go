package ud_client

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"io"
	"net/http"
	"strings"
)

type TLDRClient interface {
	ShortDefinitions(defs []string) error
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

func (c *ShortClient) ShortDefinitions(defs []string) (string, error) {

	url := "https://tldrthis.p.rapidapi.com/v1/model/abstractive/summarize-text/"
	text := strings.Join(defs, ",")
	payload := strings.NewReader(fmt.Sprintf(`{"text": "%s"}`, text))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", "3ba6283321msh80a967930e01301p13ea9fjsna2afafa73bdc")
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

	// Извлечение значения summary
	summary, ok := data["summary"].(string)
	if !ok {
		return "err", nil
	}

	return summary, nil

}
