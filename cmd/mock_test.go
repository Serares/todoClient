package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testResp = map[string]struct {
	Status int
	Body   string
}{
	"resultsMany": {
		Status: http.StatusOK,
		Body: `{
	"results": [
	{
	"Task": "Task 1",
	"Done": false,
	"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
	"CompletedAt": "0001-01-01T00:00:00Z"
	},
	{
	"Task": "Task 2",
	"Done": false,
	"CreatedAt": "2019-10-28T08:23:38.323447798-04:00",
	"CompletedAt": "0001-01-01T00:00:00Z"
	}
	],"date": 1572265440,
	"total_results": 2
	}`,
	},
	"resultsOne": {
		Status: http.StatusOK,
		Body: `{
	"results": [
	{
	"Task": "Task 1",
	"Done": false,
	"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
	"CompletedAt": "0001-01-01T00:00:00Z"
	}
	],
	"date": 1572265440,
	"total_results": 1
	}`,
	},
	"noResults": {
		Status: http.StatusOK,
		Body: `{
	"results": [],
	"date": 1572265440,
	"total_results": 0
	}`,
	},
	"root": {
		Status: http.StatusOK,
		Body:   "There's an API here",
	},
	"notFound": {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},
}

func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)

	return ts.URL, func() {
		ts.Close()
	}
}

func TestListAction(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		expOut   string
		resp     struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:     "Results",
			expError: nil,
			expOut:   "- 1 Task 1\n- 2 Task 2\n",
			resp:     testResp["resultsMany"],
		},
		{
			name:     "NoResults",
			expError: ErrNotFound,
			resp:     testResp["noResults"],
		},
		{
			name:        "InvalidURL",
			expError:    ErrConnection,
			resp:        testResp["noResults"],
			closeServer: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.resp.Status)
					fmt.Fprintln(w, tc.resp.Body)
				})
			defer cleanup()

			if tc.closeServer {
				cleanup()
			}

			var out bytes.Buffer
			err := listAction(&out, url)
			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got no error.", tc.expError)
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("Expected error %q, got %q.", tc.expError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got %q.", err)
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected output %q, got %q", tc.expOut, out.String())
			}
		})
	}
}
