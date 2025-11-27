package api

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type ShortenBatchURLRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

type ShortenBatchURLResponse struct {
	ID  string `json:"correlation_id"`
	URL string `json:"short_url"`
}

type ShortenUserURLsResponse struct {
	Alias string `json:"short_url"`
	URL   string `json:"original_url"`
}
