package ud_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type UDClient interface {
	DefineTerm(ctx context.Context, term string) ([]string, error)
}

type Client struct {
	apiKey      string
	host        string
	rateLimiter *rate.Limiter
	httpClient  *http.Client
	workerPool  chan struct{}
}

func NewUDClient(apiKey, host string, rateLimit rate.Limit, burst, maxConcurrent int) *Client {
	return &Client{
		apiKey:      apiKey,
		host:        host,
		rateLimiter: rate.NewLimiter(rateLimit, burst),
		httpClient:  &http.Client{},
		workerPool:  make(chan struct{}, maxConcurrent),
	}
}

func (c *Client) fetchDefinition(ctx context.Context, term string) ([]string, error) {
	select {
	case c.workerPool <- struct{}{}:
		defer func() { <-c.workerPool }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	url := fmt.Sprintf("https://mashape-community-urban-dictionary.p.rapidapi.com/define?term=%s", term)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-RapidAPI-Key", c.apiKey)
	req.Header.Add("X-RapidAPI-Host", c.host)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var definition UrbanDictionaryResponse
	if err := json.Unmarshal(body, &definition); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	definitions := make([]string, len(definition.List))
	for i, item := range definition.List {
		definitions[i] = item.Definition
	}

	return definitions, nil
}

func (c *Client) DefineTerm(ctx context.Context, term string) ([]string, error) {
	return c.fetchDefinition(ctx, term)
}

func (c *Client) DefineTerms(ctx context.Context, terms []string) (map[string][]string, error) {
	results := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, term := range terms {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			definitions, err := c.fetchDefinition(ctx, t)
			if err != nil {
				log.Printf("Error fetching definition for %s: %v", t, err)
				return
			}
			mu.Lock()
			results[t] = definitions
			mu.Unlock()
		}(term)
	}

	wg.Wait()
	return results, nil
}
