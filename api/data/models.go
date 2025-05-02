package data

type Cart struct {
	Cart []string `json:"cart,omitempty"`
}

type Member struct {
	Username   string   `json:"username" dynamodbav:"username"`
	FirstName  string   `json:"first_name" dynamodbav:"first_name"`
	LastName   string   `json:"last_name" dynamodbav:"last_name"`
	Cart       []string `json:"cart,omitempty" dynamodbav:"cart,omitempty"`
	Checkedout []string `json:"checked_out,omitempty" dynamodbav:"checked_out,omitempty"`
	Rented     []string `json:"rented,omitempty" dynamodbav:"rented,omitempty"`
	Type       string   `json:"member_type" dynamodbav:"member_type"`
}

type Movie struct {
	ID        string   `json:"id" dynamodbav:"id"`
	Title     string   `json:"title" dynamodbav:"title"`
	Cast      []string `json:"cast" dynamodbav:"cast"`
	Director  string   `json:"director" dynamodbav:"director"`
	Inventory int      `json:"inventory" dynamodbav:"inventory"`
	Rented    int      `json:"rented,omitempty" dynamodbav:"rented"`
	Rating    string   `json:"rating" dynamodbav:"rating"`
	Review    string   `json:"review" dynamodbav:"review"`
	Synopsis  string   `json:"synopsis" dynamodbav:"synopsis"`
	Year      string   `json:"year" dynamodbav:"year"`
}

type CartMovie struct {
	ID        string `json:"id" dynamodbav:"id"`
	Title     string `json:"title" dynamodbav:"title"`
	Invnetory int    `json:"inventory" dynamodbav:"inventory"`
}
