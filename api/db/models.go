package db

type Member struct {
	Username  string   `json:"username"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Cart      []string `json:"cart,omitempty"`
	Type      string   `json:"member_type"`
}

type Movie struct {
	Id        string
	Title     string
	Cast      []string
	Director  string
	Inventory int
	Rented    int `json:"rented,omitempty"`
	Rating    string
	Review    string
	Synopsis  string
	Year      string
}

// @TODO: add omitempty to relevant fields; add cart to member
