# API Reference

The API server listens on `localhost:8080`. REST routes use the `/api/v1` prefix. GraphQL is available at `/graphql/v1`.

## REST endpoints

### Movies

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/movies?page={page}` | List movies from a page. `page` is optional and defaults to `A`; valid page keys are `#` and `A`–`Z`. |
| `GET` | `/api/v1/movies/{movieID}` | Retrieve a movie by its ID. |
| `GET` | `/api/v1/movies/{movieID}/metrics` | Retrieve metrics for a movie. |
| `GET` | `/api/v1/movies/{movieID}/trivia` | Retrieve trivia for a movie. |

### Members and rentals

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/members/{username}` | Retrieve a member profile. |
| `POST` | `/api/v1/members/login` | Log in or create a member using a JSON body with `username`. |
| `GET` | `/api/v1/members/{username}/cart` | Retrieve the movies in a member's cart. |
| `GET` | `/api/v1/members/cart/{username}` | Alias for retrieving a member's cart. |
| `PUT` | `/api/v1/members/cart` | Add a movie to a member's cart. JSON body: `username`, `movie_id`. |
| `PUT` | `/api/v1/members/cart/remove` | Remove a movie from a member's cart. JSON body: `username`, `movie_id`. |
| `GET` | `/api/v1/members/{username}/cart/ids` | Retrieve the movie IDs in a member's cart. |
| `POST` | `/api/v1/members/checkout` | Check out movies for a member. JSON body: `username`, `movie_ids` (non-empty array). |
| `POST` | `/api/v1/members/return` | Return rented movies for a member. JSON body: `username`, `movie_ids` (non-empty array). |
| `GET` | `/api/v1/members/{username}/checkedout` | Retrieve movies currently checked out by a member. |
| `PUT` | `/api/v1/members/{username}?api_choice={choice}` | Set a member's API preference; `api_choice` must be `REST` or `GraphQL`. |

### Movie recommendations

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/members/mood/initial_voting` | Get the initial movie slate for recommendation voting. |
| `POST` | `/api/v1/members/mood/vote` | Advance a recommendation-voting iteration using the current mood, iteration, previous selection count, and movie IDs. |
| `POST` | `/api/v1/members/mood` | Update the recommendation mood using the current mood, iteration, and movie IDs. |
| `POST` | `/api/v1/members/mood/picks` | Get final recommendations for a mood; returns a best pick and additional good picks. |

## GraphQL endpoint

`POST /graphql/v1` accepts GraphQL queries and mutations in the request body.

### Queries

| Operation | Description |
|---|---|
| `GetMovies(page, director)` | List movies for a page (defaults to `A`), optionally filtered by director. |
| `GetMovie(movieID)` | Retrieve a movie by ID. |
| `GetCart(username)` | Retrieve a member's cart movies. |
| `GetCheckedout(username)` | Retrieve a member's checked-out movies. |
| `GetMember(username)` | Retrieve a member profile. |
| `DirectedMovies(director)` | Find movies directed by a director. |
| `DirectedPerformers(director)` | List performers who have worked with a director. |
| `StarredIn(star)` | Find movies featuring a performer. |
| `StarredWith(star)` | List performers who have appeared with a performer. |
| `KevinBacon(star, title, director, depth)` | Explore connected performers, movies, and directors from a performer, movie title, or director. At least one of `star`, `title`, or `director` is required; `depth` defaults to `1` and is capped at `10`. |

### Mutations

| Operation | Description |
|---|---|
| `ReturnRentals(username, movieIDs)` | Return the specified movies for a member. |
| `UpdateCart(username, movieID, removeFromCart)` | Add a movie to a member's cart, or remove it when `removeFromCart` is true. |
| `Checkout(username, movieIDs)` | Check out the specified movies for a member. |
| `SetAPIChoice(username, apiChoice)` | Set a member's API preference to `REST` or `GraphQL` (defaults to `REST`). |
