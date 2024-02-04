package ud_client

type UrbanDictionaryResponse struct {
	List []struct {
		Definition string `json:"definition"`
	} `json:"list"`
}
