package data

var TestMovies = []Movie{
	{Cast: []string{
		"Kevin Spacey",
		"Russell Crowe",
		"Guy Pearce",
		"James Cromwell",
	},
		Director:  "Curtis Hanson",
		ID:        "l.a._confidential_1997",
		Inventory: 5,
		Rented:    0,
		Rating:    "99%",
		Review:    "foo",
		Synopsis:  "bar",
		Title:     "L.A. Confidential",
		Trivia:    ":What iconic L.A. landmark is featured in the film as the setting for the bloody \"Bloody Christmas\" incident?: The Christmas-themed Central Booking lobby of the L.A. Police Department&:&:What is the name of the call girl who bears a striking resemblance to Veronica Lake and is crucial to the central plot?: Lynn Bracken&:&:Which actor won the Academy Award for Best Supporting Actress for their portrayal of Lynn Bracken in L.A. Confidential?: Kim Basinger",
		Year:      "1997",
	},
	{Cast: []string{
		"Humphrey Bogart",
		"Ingrid Bergman",
		"Paul Henreid",
		"Claude Rains",
	},
		Director:  "Michael Curtiz",
		ID:        "casablanca_1942",
		Inventory: 4,
		Rented:    0,
		Rating:    "99%",
		Review:    " An undisputed masterpiece and perhaps Hollywood's quintessential statement on love and romance, ",
		Synopsis:  "Rick Blaine (Humphrey Bogart), who owns a nightclub in Casablanca, discovers his old flame Ilsa (Ingrid Bergman) is in town...",
		Title:     "Casablanca",
		Trivia:    ":What song, frequently requested by Ilsa Lund, becomes a symbol of her past relationship with Rick and a source of conflict?:\"As Time Goes By\"&:&:What is the name of the gambling establishment owned by Rick Blaine in Casablanca?:Rick's Café Américain&:&:What are the \"letters of transit\" that everyone in Casablanca is so desperate to obtain?:They are documents that allow the bearer to travel freely to neutral Portugal and then on to the United States, effectively escaping Nazi-occupied Europe.",
		Year:      "1942",
	},
}

var TestMovieIDs = []string{
	"l.a._confidential_1997",
	"casablanca_1942",
}

var TestMember = Member{
	FirstName: "Sea",
	LastName:  "Captain",
	Username:  "sea_captain",
	Type:      "advanced",
	Cart:      TestMovieIDs,
}
