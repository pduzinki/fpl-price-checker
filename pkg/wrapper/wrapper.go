package wrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ErrHttpRequest indicates that something wrong happened in http request
type ErrHttpRequest struct {
	statusCode int
}

func (e ErrHttpRequest) Error() string {
	return fmt.Sprintf("http status not ok: %d\n", e.statusCode)
}

func (e ErrHttpRequest) GetHttpStatusCode() int {
	return e.statusCode
}

const DefaultURL = "https://fantasy.premierleague.com/api"

type Wrapper struct {
	client  *http.Client
	baseURL string
}

// NewWrapper returns new instance of Wrapper
func NewWrapper() *Wrapper {
	return &Wrapper{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: DefaultURL,
	}
}

// GetPlayers queries https://fantasy.premierleague.com/api/bootstrap-static/
// and returns slice of wrapper.Player, or error otherwise
func (w *Wrapper) GetPlayers() ([]Player, error) {
	url := fmt.Sprintf(w.baseURL + "/bootstrap-static/")
	var bs Bootstrap

	err := w.fetchData(url, &bs)
	if err != nil {
		return nil, fmt.Errorf("GetPlayers: %w", err)
	}

	return bs.Players, nil
}

// fetchData is a helper method that forms and sends http request,
// and unmarshals the response
func (w *Wrapper) fetchData(url string, data interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("fetchData, creating new request failed: %w", err)
	}
	req.Header.Set("User-Agent", "app")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("fetchData, sending http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetchData, http request failed: %w", ErrHttpRequest{resp.StatusCode})
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("fetchData, reading response body failed: %w", err)
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return fmt.Errorf("fetchData, unmarshalling data failed: %w", err)
	}

	return nil
}
