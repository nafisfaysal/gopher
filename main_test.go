package main

import (
	"net/http"
	"testing"
)

func TestParallelize(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{"StatusOk", "https://i.pinimg.com/236x/0d/a8/87/0da8872e1ca3e247aef7f75f64a75a5f--learn-coding-logos.jpg", http.StatusOK},
		{"StatusNotFound", "https://forum.golangbridge.org/uploads/default/original/2X/0/03cbc1a9f9178055093eb0c25ba9df2c29611671.png", http.StatusNotFound},
	}

	{
		for i, tt := range tests {
			tf := func(t *testing.T) {
				t.Parallel()

				t.Logf("Test: %d : checking %q for status code %d", i, tt.url, tt.statusCode)
				{
					resp, err := http.Get(tt.url)
					if err != nil {
						t.Fatalf("Should not be able to make the Get call : %v", err)
					}

					defer resp.Body.Close()

					if resp.StatusCode == tt.statusCode {
						t.Logf("Should receive a %d status code.", tt.statusCode)
					} else {
						t.Errorf("Should receive a %d status code : %v", tt.statusCode, resp.StatusCode)
					}
				}
			}

			t.Run(tt.name, tf)
		}
	}
}
