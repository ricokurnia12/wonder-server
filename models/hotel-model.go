package models

type NearestPlace struct {
	Name     string `json:"name"`
	Distance string `json:"distance"`
}

type Hotel struct {
	ID            int64          `json:"id"`
	Name          string         `json:"name"`
	Area          string         `json:"area"`
	Address       string         `json:"address"`
	PricePerNight int64          `json:"price_per_night"`
	Facilities    []string       `json:"facilities"`
	NearestPlaces []NearestPlace `json:"nearest_places"`
	ImageURL      string         `json:"image_url"`
	TravelokaURL  string         `json:"traveloka_url"`
}
