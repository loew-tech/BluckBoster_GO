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
	APIChoice  string   `json:"api_choice,omitempty" dynamodbav:"api_choice,omitempty"`
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

type MovieMetrics struct {
	Acting         float32 `json:"acting" dynamodbav:"acting"`
	Action         float32 `json:"action" dynamodbav:"action"`
	Cinematography float32 `json:"cinematography" dynamodbav:"cinematography"`
	Comedy         float32 `json:"comedy" dynamodbav:"comedy"`
	Directing      float32 `json:"directing" dynamodbav:"directing"`
	Drama          float32 `json:"drama" dynamodbav:"drama"`
	Fantasy        float32 `json:"fantasy" dynamodbav:"fantasy"`
	Horror         float32 `json:"horror" dynamodbav:"horror"`
	Romance        float32 `json:"romance" dynamodbav:"romance"`
	Suspense       float32 `json:"suspense" dynamodbav:"suspense"`
	Writing        float32 `json:"writing" dynamodbav:"writing"`
}
