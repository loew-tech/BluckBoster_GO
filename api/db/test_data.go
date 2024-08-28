package db

var TestMovies = []Movie{
	{Cast: []string{
		"Kevin Spacey",
		"Russell Crowe",
		"Guy Pearce",
		"James Cromwell",
	},
		Director:  "Curtis Hanson",
		Id:        "l.a._confidential_1997",
		Inventory: 5,
		Rented:    nil,
		Rating:    "99%",
		Review:    "foo",
		Synopsis:  "bar",
		Title:     "L.A. Confidential",
		Year:      "1997",
	},
	{Cast: []string{
		"Chris Evans",
		"Russell Crowe",
		"Guy Pearce",
		"James Cromwell",
	},
		Director:  "Russo Brothers",
		Id:        "the_winter_soldier",
		Inventory: 5,
		Rented:    nil,
		Rating:    "99%",
		Review:    "123",
		Synopsis:  "456",
		Title:     "The Winter Soldier",
		Year:      "2014",
	},
}

var TestMember = Member{
	FirstName: "Sea",
	LastName:  "Captain",
	Username:  "sea_captain",
	Type:      "advanced",
}