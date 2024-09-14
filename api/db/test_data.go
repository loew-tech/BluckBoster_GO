package db

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
