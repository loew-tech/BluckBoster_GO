package constants

const (
	ID = "id"

	// Member
	MEMBER      = "member"
	MEMBER_TYPE = "Member"
	USERNAME    = "username"
	FIRSTNAME   = "first_name"
	LASTNAME    = "last_name"
	CART_STRING = "cart"
	CHECKED_OUT = "checked_out"
	TYPE        = "member_type"

	MEMBER_TYPE_BASIC   = "basic"
	MEMBER_TYPE_ADVANCE = "advanced"
	MEMBER_TYPE_PREMIUM = "premium"

	// Movie
	MOVIE      = "movie"
	MOVIES     = "movies"
	MOVIE_TYPE = "Movie"
	TITLE      = "title"
	CAST       = "cast"
	YEAR       = "year"
	TRIVIA     = "trivia"
	INVENTORY  = "inventory"
	DIRECTOR   = "director"
	RATING     = "rating"
	RENTED     = "rented"
	SYNOPSIS   = "synopsis"
	REVIEW     = "review"

	// Cart
	ADD    = "ADD"
	DELETE = "DELETE"

	CART         = true
	NOT_CART     = false
	CHECKOUT     = true
	NOT_CHECKOUT = false

	// GraphQL
	MOVIE_ID    = "movieID"
	PAGE        = "page"
	GET_MOVIES  = "GetMovies"
	GET_MOVIE   = "GetMovie"
	GET_MEMBER  = "GetMember"
	DIRECTED_BY = "DirectedBy"
	STAR        = "star"
	STARREDIN   = "StarredIn"
	STARREDWITH = "StarredWith"

	// AWS
	PAGINATE_KEY       = "paginate_key"
	PAGINATE_KEY_INDEX = "paginate_key-index"
)
