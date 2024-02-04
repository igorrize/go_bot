package ud_client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type UDClient interface {
	DefineTerm(term string) error
}

type Client struct {
	apiKey string
	host   string
}

func NewUDClient(apiKey, host string) *Client {
	return &Client{
		apiKey: apiKey,
		host:   host,
	}
}

func (c *Client) DefineTerm(term string) ([]string, error) {
	url := fmt.Sprintf("https://mashape-community-urban-dictionary.p.rapidapi.com/define?term=%s", term)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", c.apiKey)
	req.Header.Add("X-RapidAPI-Host", c.host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var definition UrbanDictionaryResponse
	jsonErr := json.Unmarshal(body, &definition)
	if jsonErr != nil {
		return nil, err
	}
    var allDefinitions []string

    for _, item := range definition.List {
		allDefinitions = append(allDefinitions, item.Definition)
	}

	return allDefinitions, nil
	
}
