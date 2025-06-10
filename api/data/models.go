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
	Cast      []string `json:"cast,omitempty" dynamodbav:"cast,omitempty"`
	Director  string   `json:"director,omitempty" dynamodbav:"director,omitempty"`
	Inventory int      `json:"inventory,omitempty" dynamodbav:"inventory,omitempty"`
	Rented    int      `json:"rented,omitempty" dynamodbav:"rented,omitempty"`
	Rating    string   `json:"rating,omitempty" dynamodbav:"rating,omitempty"`
	Review    string   `json:"review,omitempty" dynamodbav:"review,omitempty"`
	Synopsis  string   `json:"synopsis,omitempty" dynamodbav:"synopsis,omitempty"`
	Trivia    string   `json:"trivia,omitempty" dynamodbav:"trivia,omitempty"`
	Year      string   `json:"year,omitempty" dynamodbav:"year,omitempty"`
}

type MovieTrivia struct {
	Trivia string `json:"trivia" dynamodbav:"trivia"`
}
