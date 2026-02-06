package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const (
	DefaultUserAgent = "ManniWatch Api Client/Go"
)

type ManniwatchClient struct {
	httpClient    *http.Client
	host          string
	settingsMutex sync.RWMutex
}

func NewManniwatchClient(host string) *ManniwatchClient {
	return &ManniwatchClient{
		httpClient: &http.Client{},
		host:       host,
	}
}

func (c *ManniwatchClient) SetHost(host string) {
	c.settingsMutex.Lock()
	defer c.settingsMutex.Unlock()
	c.host = host
}

func (c *ManniwatchClient) GetHost() string {
	c.settingsMutex.RLock()
	defer c.settingsMutex.RUnlock()
	return c.host
}

func (c *ManniwatchClient) buildURL(path string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   c.GetHost(),
		Path:   path,
	}
}

// buildRequest helper handles method, path, and optional query params or body
func (c *ManniwatchClient) buildRequest(method, path string, queryParams url.Values, body io.Reader) (*http.Request, error) {
	u := c.buildURL(path)
	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", DefaultUserAgent)

	// Default to JSON accept
	req.Header.Set("Accept", "application/json")

	// Set Content-Type if body is present (assuming form-urlencoded usually based on TS qs.stringify usage)
	// If body is nil, no Content-Type
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

// executeRequest performs the request and decodes the JSON response into v
func (c *ManniwatchClient) executeRequest(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	if v == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// GetVehicleLocations returns the vehicle locations
// positionType: coordinate type (default: CORRECTED)
// lastUpdate: timestamp of last update
func (c *ManniwatchClient) GetVehicleLocations(positionType PositionType, lastUpdate int64) (*VehicleLocationList, error) {
	if positionType == "" {
		positionType = PositionTypeCorrected
	}

	params := url.Values{}
	params.Set("colorType", "ROUTE_BASED")
	params.Set("positionType", string(positionType))
	if lastUpdate > 0 {
		params.Set("lastUpdate", strconv.FormatInt(lastUpdate, 10))
	}

	req, err := c.buildRequest(http.MethodGet, "/internetservice/geoserviceDispatcher/services/vehicleinfo/vehicles", params, nil)
	if err != nil {
		return nil, err
	}

	var result VehicleLocationList
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRouteByTripId returns route info for a trip
func (c *ManniwatchClient) GetRouteByTripId(tripID string) (*VehiclePathInfo, error) {
	params := url.Values{}
	params.Set("id", tripID)

	// Note: TS uses POST with params (query string) for this endpoint
	req, err := c.buildRequest(http.MethodPost, "/internetservice/geoserviceDispatcher/services/pathinfo/trip", params, nil)
	if err != nil {
		return nil, err
	}

	var result VehiclePathInfo
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRouteByVehicleId returns route info for a vehicle
func (c *ManniwatchClient) GetRouteByVehicleId(vehicleID string) (*VehiclePathInfo, error) {
	params := url.Values{}
	params.Set("id", vehicleID)

	req, err := c.buildRequest(http.MethodPost, "/internetservice/geoserviceDispatcher/services/pathinfo/vehicle", params, nil)
	if err != nil {
		return nil, err
	}

	var result VehiclePathInfo
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRouteByRouteId returns route info for a route
func (c *ManniwatchClient) GetRouteByRouteId(routeID string, direction string) (*VehiclePathInfo, error) {
	params := url.Values{}
	params.Set("id", routeID)
	params.Set("direction", direction)

	req, err := c.buildRequest(http.MethodPost, "/internetservice/geoserviceDispatcher/services/pathinfo/route", params, nil)
	if err != nil {
		return nil, err
	}

	var result VehiclePathInfo
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStopLocations returns stop locations within a bounding box
func (c *ManniwatchClient) GetStopLocations(box BoundingBox) (*StopLocations, error) {
	params := url.Values{}
	params.Set("top", strconv.FormatInt(box.Top, 10))
	params.Set("bottom", strconv.FormatInt(box.Bottom, 10))
	params.Set("left", strconv.FormatInt(box.Left, 10))
	params.Set("right", strconv.FormatInt(box.Right, 10))

	req, err := c.buildRequest(http.MethodGet, "/internetservice/geoserviceDispatcher/services/stopinfo/stops", params, nil)
	if err != nil {
		return nil, err
	}

	var result StopLocations
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStopPointLocations returns stop point locations within a bounding box
func (c *ManniwatchClient) GetStopPointLocations(box BoundingBox) (*StopPointLocations, error) {
	params := url.Values{}
	params.Set("top", strconv.FormatInt(box.Top, 10))
	params.Set("bottom", strconv.FormatInt(box.Bottom, 10))
	params.Set("left", strconv.FormatInt(box.Left, 10))
	params.Set("right", strconv.FormatInt(box.Right, 10))

	req, err := c.buildRequest(http.MethodGet, "/internetservice/geoserviceDispatcher/services/stopinfo/stopPoints", params, nil)
	if err != nil {
		return nil, err
	}

	var result StopPointLocations
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTripPassages returns trip passages
func (c *ManniwatchClient) GetTripPassages(tripID string, mode StopMode) (*TripPassages, error) {
	if mode == "" {
		mode = StopModeDeparture
	}

	data := url.Values{}
	data.Set("tripId", tripID)
	data.Set("mode", string(mode))

	// TS uses data (body) for this endpoint
	req, err := c.buildRequest(http.MethodPost, "/internetservice/services/tripInfo/tripPassages", nil, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	var result TripPassages
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStopPassages returns stop passages
func (c *ManniwatchClient) GetStopPassages(stopID string, mode StopMode, startTime int64, timeFrame int) (*StopPassage, error) {
	if mode == "" {
		mode = StopModeDeparture
	}

	data := url.Values{}
	data.Set("stop", stopID)
	data.Set("mode", string(mode))
	if startTime > 0 {
		data.Set("startTime", strconv.FormatInt(startTime, 10))
	}
	if timeFrame > 0 {
		data.Set("timeFrame", strconv.Itoa(timeFrame))
	}

	req, err := c.buildRequest(http.MethodPost, "/internetservice/services/passageInfo/stopPassages/stop", nil, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	var result StopPassage
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStopPointPassages returns stop point passages
func (c *ManniwatchClient) GetStopPointPassages(stopPointID string, mode StopMode, startTime int64, timeFrame int) (*StopPassage, error) {
	if mode == "" {
		mode = StopModeDeparture
	}

	data := url.Values{}
	data.Set("stopPoint", stopPointID)
	data.Set("mode", string(mode))
	if startTime > 0 {
		data.Set("startTime", strconv.FormatInt(startTime, 10))
	}
	if timeFrame > 0 {
		data.Set("timeFrame", strconv.Itoa(timeFrame))
	}

	req, err := c.buildRequest(http.MethodPost, "/internetservice/services/passageInfo/stopPassages/stopPoint", nil, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	var result StopPassage
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStopInfo returns stop info
func (c *ManniwatchClient) GetStopInfo(stopID string) (*StopInfo, error) {
	data := url.Values{}
	data.Set("stop", stopID)

	req, err := c.buildRequest(http.MethodPost, "/internetservice/services/stopInfo/stop", nil, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	var result StopInfo
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStopPointInfo returns stop point info
func (c *ManniwatchClient) GetStopPointInfo(stopPointID string) (*StopPointInfo, error) {
	data := url.Values{}
	data.Set("stopPoint", stopPointID)

	req, err := c.buildRequest(http.MethodPost, "/internetservice/services/stopInfo/stopPoint", nil, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	var result StopPointInfo
	err = c.executeRequest(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSettings returns settings
func (c *ManniwatchClient) GetSettings() (*Settings, error) {
	req, err := c.buildRequest(http.MethodGet, "/internetservice/settings", nil, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/javascript")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(bodyBytes)

	bracketStart := strings.Index(bodyStr, "{")
	bracketEnd := strings.LastIndex(bodyStr, "}")

	if bracketStart >= 0 && bracketEnd > bracketStart {
		jsonStr := bodyStr[bracketStart : bracketEnd+1]
		var result Settings
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			return nil, err
		}
		return &result, nil
	}

	return nil, fmt.Errorf("non valid response body")
}
