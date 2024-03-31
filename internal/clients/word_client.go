package ud_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WordClient interface {
	WordDefinition(term string) error
}

type WordApiClient struct {
	apiKey string
	host   string
}

func NewWordApiClient(apiKey, host string) *WordApiClient {
	return &WordApiClient{
		apiKey: apiKey,
		host: host,
		}
}

func (c *WordApiClient) WordDefinition(term string) string {
	url := fmt.Sprintf("https://wordsapiv1.p.rapidapi.com/words/%s/definitions", term)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", c.apiKey)
	req.Header.Add("X-RapidAPI-Host", c.host)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var definitionResponse WordApiResponse

	err := json.Unmarshal(body, &definitionResponse)
	if err != nil {
		fmt.Println("err JSON:", err)
		return "err JSON"
	}
	
	fmt.Println(definitionResponse)

	return definitionResponse.Definitions[0].Definition

}

type WordApiResponse struct {
	Word        string `json:"word"`
	Definitions []struct {
		Definition   string `json:"definition"`
		PartOfSpeech string `json:"partOfSpeech"`
	} `json:"definitions"`
}
