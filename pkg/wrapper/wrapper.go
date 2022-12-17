package wrapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const DefaultURL = "https://fantasy.premierleague.com/api"

type Wrapper struct {
	client  *http.Client
	baseURL string
}

// NewWrapper returns new instance of Wrapper
func NewWrapper() Wrapper {
	return Wrapper{
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
		return nil, err
	}

	return bs.Players, nil
}

// fetchData is a helper method that forms and sends http request,
// and unmarshals the response
func (w *Wrapper) fetchData(url string, data interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "app")

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO consider improving/adding custom error types

	if resp.StatusCode != http.StatusOK {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	return nil
}
