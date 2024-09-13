package db

type Cart struct {
	Cart []string `json:"cart,omitempty"`
}

type Member struct {
	Username  string   `json:"username"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Cart      []string `json:"cart,omitempty"`
	Type      string   `json:"member_type"`
}

type Movie struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Cast      []string `json:"cast"`
	Director  string   `json:"director"`
	Inventory int      `json:"inventory"`
	Rented    int      `json:"rented,omitempty"`
	Rating    string   `json:"rating"`
	Review    string   `json:"review"`
	Synopsis  string   `json:"synopsis"`
	Year      string   `json:"year"`
}

type CartMovie struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Inventory int    `json:"inventory"`
}
