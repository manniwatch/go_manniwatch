package client

// StopInfo represents basic stop passenger information
type StopInfo struct {
	ID            string `json:"id"`
	PassengerName string `json:"passengerName"`
}

// VehicleCategory represents the type of vehicle (bus or tram)
type VehicleCategory string

const (
	VehicleCategoryBus  VehicleCategory = "bus"
	VehicleCategoryTram VehicleCategory = "tram"
)

// StopLocation Information about the stop location
// Category: geo
type StopLocation struct {
	// Type of vehicle
	Category VehicleCategory `json:"category"`
	// Stop Id
	ID string `json:"id"`
	// Latitude in arcmiliseconds
	Latitude int64 `json:"latitude"`
	// Longitutde in arcmiliseconds
	Longitude int64 `json:"longitude"`
	// Humanreadable name of the stop
	Name string `json:"name"`
	// Stop short id
	ShortName string `json:"shortName"`
}

// StopLocations Response for /internetservice/geoserviceDispatcher/services/stopinfo/stops
type StopLocations struct {
	Stops []StopLocation `json:"stops"`
}

// StopMode to query
type StopMode string

const (
	StopModeArrival   StopMode = "arrival"
	StopModeDeparture StopMode = "departure"
)

// PositionType Should coordinates be corrected or unaltered
type PositionType string

const (
	PositionTypeCorrected PositionType = "CORRECTED"
	PositionTypeRaw       PositionType = "RAW"
)

// VehicleStatus Vehicle status
type VehicleStatus string

const (
	VehicleStatusPredicted VehicleStatus = "PREDICTED"
	VehicleStatusDeparted  VehicleStatus = "DEPARTED"
	VehicleStatusStopping  VehicleStatus = "STOPPING"
	VehicleStatusPlanned   VehicleStatus = "PLANNED"
)

// Departure information for vehicles
type Departure struct {
	// Time in seconds estimated from server time to arrival
	ActualRelativeTime int64 `json:"actualRelativeTime"`
	// Time of arrivabl formated HH:mm
	ActualTime  string `json:"actualTime"`
	Direction   string `json:"direction"`
	MixedTime   string `json:"mixedTime"`
	PassageID   string `json:"passageid"`
	PatternText string `json:"patternText"`
	// Planned time of arrival
	PlannedTime string `json:"plannedTime"`
	RouteID     string `json:"routeId"`
	// Current status of the vehicle
	Status    VehicleStatus `json:"status"`
	TripID    string        `json:"tripId"`
	VehicleID string        `json:"vehicleId"`
}

// RouteAlert Alerts on route
type RouteAlert struct {
	Direction   []string `json:"direction"`
	DirectionID string   `json:"directionId"`
	Title       string   `json:"title"`
}

// Route of the vehicle
type Route struct {
	Alerts     []RouteAlert    `json:"alerts"`
	Authority  string          `json:"authority"`
	Directions []string        `json:"directions"`
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	RouteType  VehicleCategory `json:"routeType"`
	// route short name
	ShortName string `json:"shortName"`
}

// StopPassage
type StopPassage struct {
	// Actual/Future Departures
	Actual           []Departure   `json:"actual"`
	Directions       []interface{} `json:"directions"` // directions: unknown[]
	FirstPassageTime int64         `json:"firstPassageTime"`
	// TODO: Need schema
	GeneralAlerts   []interface{} `json:"generalAlerts"` // generalAlerts: unknown[]
	LastPassageTime int64         `json:"lastPassageTime"`
	// Previous departures
	Old []Departure `json:"old"`
	// Routes served by this stop
	Routes []Route `json:"routes"`
	// Human readable string of the stop
	StopName string `json:"stopName"`
	// short name of the stop
	StopShortName string `json:"stopShortName"`
}

// PathSegment represents a segment of the path
type PathSegment struct {
	// Start y coordinate in arcmiliseconds
	Y1 int64 `json:"y1"`
	// Start x coordinate in arcmiliseconds
	X1 int64 `json:"x1"`
	// End y coordinate in arcmiliseconds
	Y2 int64 `json:"y2"`
	// End x coordinate in arcmiliseconds
	X2 int64 `json:"x2"`
	// Angle in degrees
	Angle float64 `json:"angle"`
	// Distance between start and end
	Length int64 `json:"length"`
}

// VehicleLocation represents a vehicle's position and status
// It handles both active and deleted vehicles via IsDeleted field
type VehicleLocation struct {
	// Vehicle Id
	ID string `json:"id"`
	// IsDeleted is true if data is stale/non-existent
	IsDeleted bool `json:"isDeleted,omitempty"`

	// --- Fields for IVehicleLocation ---

	// Kind of Vehicle
	Category VehicleCategory `json:"category,omitempty"`
	Color    string          `json:"color,omitempty"`
	// Heading of the vehicle in degrees
	Heading float64 `json:"heading,omitempty"`
	// Latitude in arcmiliseconds
	Latitude int64 `json:"latitude,omitempty"`
	// Longitude in arcmiliseconds
	Longitude int64 `json:"longitude,omitempty"`
	// Humanreadable vehicle name
	Name string `json:"name,omitempty"`
	// Previous Vehicle locations
	Path []PathSegment `json:"path,omitempty"`
	// Current TripId of the vehicle
	TripID string `json:"tripId,omitempty"`

	// --- Fields for TimestampVehicleLocation ---

	// lastUpdate attribute copied from server response
	LastUpdate int64 `json:"lastUpdate,omitempty"`
}

// VehicleLocationList Response containing list of vehicle locations
type VehicleLocationList struct {
	// Timestamp
	LastUpdate int64 `json:"lastUpdate"`
	// reported locations
	Vehicles []VehicleLocation `json:"vehicles"`
}

// BoundingBox represents a geographic bounding box
type BoundingBox struct {
	Top    int64 `json:"top"`
	Bottom int64 `json:"bottom"`
	Left   int64 `json:"left"`
	Right  int64 `json:"right"`
}

// WayPoint represents a vehicle path waypoint
type WayPoint struct {
	// Latitude in arcms
	Lat int64 `json:"lat"`
	// Longitude in arcms
	Lon int64 `json:"lon"`
	// Sequence Number
	Seq string `json:"seq"`
}

// VehiclePath represents a visual path on the map
type VehiclePath struct {
	// Color to be used to represent the path
	Color string `json:"color"`
	// Waypoints of Vehicle
	WayPoints []WayPoint `json:"wayPoints"`
}

// VehiclePathInfo Previous path points
type VehiclePathInfo struct {
	// Vehicle Paths Information
	Paths []VehiclePath `json:"paths"`
}

// StopPointLocation Information about the stop location
// Category: geo
type StopPointLocation struct {
	// Type of vehicle
	Category VehicleCategory `json:"category"`
	// Stop Id
	ID string `json:"id"`
	// label
	Label string `json:"label"`
	// Latitude in arcmiliseconds
	Latitude int64 `json:"latitude"`
	// Longitutde in arcmiliseconds
	Longitude int64 `json:"longitude"`
	// Humanreadable name of the stop
	Name string `json:"name"`
	// Stop short id
	ShortName string `json:"shortName"`
	// Stop point
	StopPoint string `json:"stopPoint"`
}

// StopPointLocations Response for /internetservice/geoserviceDispatcher/services/stopinfo/stopPoints
type StopPointLocations struct {
	// List of StopLocation
	StopPoints []StopPointLocation `json:"stopPoints"`
}

// TripPassage represents a stop on a trip
type TripPassage struct {
	// the actual estimated time
	ActualTime string `json:"actualTime,omitempty"`
	// the planned time
	PlannedTime string        `json:"plannedTime,omitempty"`
	Status      VehicleStatus `json:"status"`
	Stop        struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
	} `json:"stop"`
	StopSeqNum string `json:"stop_seq_num"`
}

// TripPassages Response from /internetservice/services/tripInfo/tripPassages
type TripPassages struct {
	// Current/Future passages
	Actual []TripPassage `json:"actual"`
	// Previous departures (Old is effectively same structure as Actual but status is DEPARTED)
	Old []TripPassage `json:"old"`
	// Directiontext
	DirectionText string `json:"directionText"`
	// Route name
	RouteName string `json:"routeName"`
}

// StopPointInfo Stop Point Information
type StopPointInfo struct {
	ID            string `json:"id"`
	PassengerName string `json:"passengerName"`
	StopPointCode string `json:"stopPointCode"`
}

type Settings struct {
	// eg. de en ar
	AvailableLanguages []string `json:"AVAILABLE_LANGUAGES"`
	DefaultTimePreview int      `json:"DEFAULT_TIME_PREVIEW"`
	GeolocationEnabled bool     `json:"GEOLOCATION_ENABLED"`
	InitialLat         float64  `json:"INITIAL_LAT"`
	InitialLon         float64  `json:"INITIAL_LON"`
	InitialZoom        int      `json:"INITIAL_ZOOM"`
	// language code
	// eg. de en ar
	Language                       string `json:"LANGUAGE"`
	MapEnabled                     bool   `json:"MAP_ENABLED"`
	MapShowControls                bool   `json:"MAP_SHOW_CONTROLS"`
	MapShowPatterns                bool   `json:"MAP_SHOW_PATTERNS"`
	MapShowStops                   bool   `json:"MAP_SHOW_STOPS"`
	MapShowVehicles                bool   `json:"MAP_SHOW_VEHICLES"`
	MaxZoom                        int    `json:"MAX_ZOOM"`
	MinZoom                        int    `json:"MIN_ZOOM"`
	MobileEnabled                  bool   `json:"MOBILE_ENABLED"`
	SearchByRoutesEnabled          bool   `json:"SEARCH_BY_ROUTES_ENABLED"`
	SearchByStoppointsEnabled      bool   `json:"SEARCH_BY_STOPPOINTS_ENABLED"`
	ShowAboutDepartureText         bool   `json:"SHOW_ABOUT_DEPARTURE_TEXT"`
	ShowActualColumn               bool   `json:"SHOW_ACTUAL_COLUMN"`
	ShowDepartingText              bool   `json:"SHOW_DEPARTING_TEXT"`
	ShowDepArrText                 bool   `json:"SHOW_DEP_ARR_TEXT"`
	ShowLanguageBar                bool   `json:"SHOW_LANGUAGE_BAR"`
	ShowMixedColumn                bool   `json:"SHOW_MIXED_COLUMN"`
	ShowPassagetypeColumn          bool   `json:"SHOW_PASSAGETYPE_COLUMN"`
	ShowScheduleColumn             bool   `json:"SHOW_SCHEDULE_COLUMN"`
	SuppressCountdownTimeIncrement bool   `json:"SUPPRESS_COUNTDOWN_TIME_INCREMENT"`
	TimesliderEnabled              bool   `json:"TIMESLIDER_ENABLED"`
}
