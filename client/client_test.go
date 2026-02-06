package client

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// roundTripFunc is a type that implements http.RoundTripper
type roundTripFunc func(req *http.Request) *http.Response

// RoundTrip executes the mock function
func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// newTestClient returns a *http.Client with Transport replaced to avoid network calls
func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

func TestGetVehicleLocations(t *testing.T) {
	mockResponse := `{
		"lastUpdate": 12345678,
		"vehicles": [
			{"id": "v1", "latitude": 100, "longitude": 200, "category": "bus"},
			{"id": "v2", "isDeleted": true}
		]
	}`

	client := NewManniwatchClient("api.example.com")
	client.httpClient = newTestClient(func(req *http.Request) *http.Response {
		expectedPath := "/internetservice/geoserviceDispatcher/services/vehicleinfo/vehicles"
		if req.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, req.URL.Path)
		}

		q := req.URL.Query()
		if q.Get("colorType") != "ROUTE_BASED" {
			t.Errorf("Expected colorType=ROUTE_BASED, got %s", q.Get("colorType"))
		}
		if q.Get("positionType") != "CORRECTED" {
			t.Errorf("Expected positionType=CORRECTED, got %s", q.Get("positionType"))
		}
		if q.Get("lastUpdate") != "1000" {
			t.Errorf("Expected lastUpdate=1000, got %s", q.Get("lastUpdate"))
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			Header:     make(http.Header),
		}
	})

	vehicles, err := client.GetVehicleLocations(PositionTypeCorrected, 1000)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if vehicles.LastUpdate != 12345678 {
		t.Errorf("Expected lastUpdate 12345678, got %d", vehicles.LastUpdate)
	}
	if len(vehicles.Vehicles) != 2 {
		t.Errorf("Expected 2 vehicles, got %d", len(vehicles.Vehicles))
	}
	if vehicles.Vehicles[0].ID != "v1" || vehicles.Vehicles[0].Category != "bus" {
		t.Errorf("Vehicle 1 data mismatch")
	}
	if !vehicles.Vehicles[1].IsDeleted {
		t.Errorf("Vehicle 2 should be deleted")
	}
}

func TestGetStopPassages(t *testing.T) {
	mockResponse := `{
		"stopName": "Test Stop",
		"actual": []
	}`

	client := NewManniwatchClient("api.example.com")
	client.httpClient = newTestClient(func(req *http.Request) *http.Response {
		expectedPath := "/internetservice/services/passageInfo/stopPassages/stop"
		if req.URL.Path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", req.Method)
		}

		// Reading body consumes it, so we need to restore it usually, but this is a one-off mock
		bodyBytes, _ := io.ReadAll(req.Body)
		bodyStr := string(bodyBytes)

		// Check for URL encoded body
		if !strings.Contains(bodyStr, "stop=test-stop-id") {
			t.Errorf("Body missing stop param: %s", bodyStr)
		}
		if !strings.Contains(bodyStr, "mode=departure") {
			t.Errorf("Body missing mode param: %s", bodyStr)
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			Header:     make(http.Header),
		}
	})

	passages, err := client.GetStopPassages("test-stop-id", StopModeDeparture, 0, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if passages.StopName != "Test Stop" {
		t.Errorf("Expected stop name 'Test Stop', got '%s'", passages.StopName)
	}
}

func TestGetSettings(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		expectError  bool
		expectedLang string
	}{
		{
			name:         "Clean JSON",
			responseBody: `{"LANGUAGE": "en"}`,
			expectError:  false,
			expectedLang: "en",
		},
		{
			name:         "Wrapped in JS",
			responseBody: `var settings = {"LANGUAGE": "de"};`,
			expectError:  false,
			expectedLang: "de",
		},
		{
			name:         "Complex JS",
			responseBody: `abc {  "LANGUAGE": "ar" } xyz`,
			expectError:  false,
			expectedLang: "ar",
		},
		{
			name:         "Invalid",
			responseBody: `no brackets here`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewManniwatchClient("api.example.com")
			client.httpClient = newTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(tt.responseBody)),
					Header:     make(http.Header),
				}
			})

			settings, err := client.GetSettings()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if settings.Language != tt.expectedLang {
				t.Errorf("Expected language '%s', got '%s'", tt.expectedLang, settings.Language)
			}
		})
	}
}
