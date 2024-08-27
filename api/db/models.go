package db

type Member struct {
	Username  string
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"member_type"`
}

type Movie struct {
	Id        string
	Title     string
	Cast      []string
	Director  string
	Inventory int
	Rented    *int
	Rating    string
	Review    string
	Synopsis  string
	Year      string
}