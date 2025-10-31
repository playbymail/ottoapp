// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package rentals implements test data for EmberJS Super Rentals.
package rentals

type Store struct {
	rentals []Rental
}

func New() (*Store, error) {
	return &Store{
		rentals: []Rental{
			{
				Type: "rental",
				Id:   "grand-old-mansion",
				Attributes: Attributes{
					Title: "Grand Old Mansion",
					Owner: "Veruca Salt",
					City:  "San Francisco",
					Location: Location{
						Lat: 37.7749,
						Lng: -122.4194,
					},
					Category:    "Estate",
					Bedrooms:    15,
					ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/c/cb/Crane_estate_(5).jpg",
					Description: "This grand old mansion sits on over 100 acres of rolling hills and dense redwood forests.",
				},
			},
			{
				Type: "rental",
				Id:   "urban-living",
				Attributes: Attributes{
					Title: "Urban Living",
					Owner: "Mike Teavee",
					City:  "Seattle",
					Location: Location{
						Lat: 47.6062,
						Lng: -122.3321,
					},
					Category:    "Condo",
					Bedrooms:    1,
					ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/2/20/Seattle_-_Barnes_and_Bell_Buildings.jpg",
					Description: "A commuters dream. This rental is within walking distance of 2 bus stops and the Metro.",
				},
			},
			{
				Type: "rental",
				Id:   "downtown-charm",
				Attributes: Attributes{
					Title: "Downtown Charm",
					Owner: "Violet Beauregarde",
					City:  "Portland",
					Location: Location{
						Lat: 45.5175,
						Lng: -122.6801,
					},
					Category:    "Apartment",
					Bedrooms:    3,
					ImageURL:    "https://upload.wikimedia.org/wikipedia/commons/f/f7/Wheeldon_Apartment_Building_-_Portland_Oregon.jpg",
					Description: "Convenience is at your doorstep with this charming downtown rental. Great restaurants and active night life are within a few feet.",
				},
			},
		},
	}, nil
}

func (s *Store) FetchRental(id string) (Rental, bool) {
	for _, r := range s.rentals {
		if r.Id == id {
			return r, true
		}
	}
	return Rental{}, false
}

func (s *Store) FetchRentals() ([]Rental, bool) {
	return s.rentals, len(s.rentals) != 0
}

type Rental struct {
	Type       string     `json:"type,omitempty"`
	Id         string     `json:"id,omitempty"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Title       string   `json:"title,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	City        string   `json:"city,omitempty"`
	Location    Location `json:"location"`
	Category    string   `json:"category,omitempty"`
	Bedrooms    int      `json:"bedrooms,omitempty"`
	ImageURL    string   `json:"image,omitempty"`
	Description string   `json:"description,omitempty"`
}
type Location struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}
