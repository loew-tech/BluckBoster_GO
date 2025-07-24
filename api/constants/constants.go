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
	MOVIE            = "movie"
	MOVIES           = "movies"
	MOVIE_TYPE       = "Movie"
	TITLE            = "title"
	CAST             = "cast"
	YEAR             = "year"
	TRIVIA           = "trivia"
	INVENTORY        = "inventory"
	DIRECTOR         = "director"
	RATING           = "rating"
	RENTED           = "rented"
	SYNOPSIS         = "synopsis"
	REVIEW           = "review"
	RETURN_MOVIE_INC = 1
	RENT_MOVIE_INC   = -1

	// Cart
	ADD    = "ADD"
	DELETE = "DELETE"

	CART         = true
	NOT_CART     = false
	CHECKOUT     = true
	NOT_CHECKOUT = false

	// GraphQL
	ROOT_QUERY          = "RootQuery"
	ROOT_MUTATION       = "RootMutation"
	CHECKOUT_STRING     = "Checkout"
	MOVIE_ID            = "movieID"
	MOVIE_IDS           = "movieIDs"
	GET_MOVIES          = "GetMovies"
	GET_MOVIE           = "GetMovie"
	GET_CHECKEDOUT      = "GetCheckedout"
	GET_CART            = "GetCart"
	RETURN_RENTALS      = "ReturnRentals"
	GET_MEMBER          = "GetMember"
	DIRECTED_MOVIES     = "DirectedMovies"
	DIRECTED_PERFORMERS = "DirectedPerformers"
	DIRECTORS           = "directors"
	STAR                = "star"
	STARS               = "stars"
	STARREDIN           = "StarredIn"
	STARREDWITH         = "StarredWith"
	KEVING_BACON        = "KevinBacon"
	KEVING_BACON_TYPE   = "KevinBaconType"
	DEPTH               = "depth"
	TOTAL_DIRECTORS     = "total_directors"
	TOTAL_MOVIES        = "total_movies"
	TOTAL_STARS         = "total_stars"
	REMOVE_FROM_CART    = "removeFromCart"
	UPDATE_CART         = "UpdateCart"
	SET_API_CHOICE      = "SetAPIChoice"
	SUCCESS             = "success"
	CODE                = "code"
	FOR_GRAPH           = true
	NOT_FOR_GRAPH       = false

	// AWS
	PAGE               = "page"
	DEFAULT_PAGE       = "A"
	PAGINATE_KEY       = "paginate_key"
	PAGINATE_KEY_INDEX = "paginate_key-index"
	PAGES              = "#ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// API Choice
	API_CHOICE  = "api_choice"
	REST_API    = "REST"
	GRAPHQL_API = "GraphQL"

	// Router group
	REST_ROUTER_GROUP = "/api/v1"
	GRAPHQL_ENDPOINT  = "/graphql/v1"
)
