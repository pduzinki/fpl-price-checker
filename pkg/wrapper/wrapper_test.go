package wrapper

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFetchData(t *testing.T) {
	testcases := []struct {
		name         string
		statusCode   int
		bodyFilePath string
		expectErr    bool
		wantErr      error
	}{
		{
			name:         "sunny scenario",
			statusCode:   http.StatusOK,
			bodyFilePath: "./testdata/fetchdata.json",
			wantErr:      nil,
		},
		{
			name:         "too many requests",
			statusCode:   http.StatusTooManyRequests,
			bodyFilePath: "./testdata/fetchdata.json",
			wantErr:      ErrHttpRequest{http.StatusTooManyRequests},
		},
	}

	for _, test := range testcases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.statusCode)
			w.Header().Set("Content-Type", "application/json")

			f, err := os.ReadFile(test.bodyFilePath)
			if err != nil {
				t.Error(err)
			}

			w.Write(f)
		}))
		defer server.Close()

		w := Wrapper{
			client:  &http.Client{},
			baseURL: server.URL,
		}

		type tmp struct {
			Data int `json:"data"`
		}
		var data tmp

		err := w.fetchData(w.baseURL, &data)
		if err != test.wantErr {
			t.Errorf("want: %d, got: %d", test.wantErr, err)
		}
	}
}
