package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	repos "blockbuster/api/repos"
)

var memberRepo = repos.NewMemberRepoWithDynamo()

func GetMemberEndpoint(c *gin.Context) {
	member, err := memberRepo.GetMemberByUsername(c, c.Param("username"), constants.NOT_CART)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve user"})
	} else {
		if member.Username != "" {
			c.IndentedJSON(http.StatusOK, member)
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("Failed to find user %s", c.Param("username"))})
		}
	}
}

type UsernameReq struct {
	Username string `json:"username"`
}

func MemberLoginEndpoint(c *gin.Context) {
	un := UsernameReq{}
	err := c.BindJSON(&un)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"msg": "Failed to unmarshall data into request"},
		)
		return
	}
	member, err := memberRepo.GetMemberByUsername(c, un.Username, constants.NOT_CART)
	if err != nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"msg": "Error retrieving user"},
		)
		return
	}
	if member.Username != "" {
		c.IndentedJSON(http.StatusOK, member)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("Failed to find user %s", c.Param("username"))})
	}
}

func GetCartIDsEndpoint(c *gin.Context) {
	user, err := memberRepo.GetMemberByUsername(c, c.Param("username"), constants.CART)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve user cart"})
	} else {
		c.IndentedJSON(http.StatusOK, user.Cart)
	}
}

func GetCartMoviesEndpoint(c *gin.Context) {
	movies, err := memberRepo.GetCartMovies(c, c.Param("username"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve cart movies"})
	} else {
		c.IndentedJSON(http.StatusOK, movies)
	}
}

type ModifyCartRequest struct {
	Username string `json:"username"`
	MovieID  string `json:"movie_id"`
}

func AddToCartEndpoint(c *gin.Context) {
	cartHelper(c, constants.ADD, false)
}

func RemoveFromCartEndpoint(c *gin.Context) {
	cartHelper(c, constants.DELETE, false)
}

func cartHelper(c *gin.Context, action string, checkingOut bool) {
	req := ModifyCartRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"msg": "Bad Request for ModifyCart"},
		)
		return
	}

	inserted, response, err := memberRepo.ModifyCart(c, req.Username, req.MovieID, action, checkingOut)
	if err != nil {
		act, direction := "adding", "to"
		if action == constants.DELETE {
			act, direction = "removing", "from"
		}
		msg := fmt.Sprintf("Error %s %s %s %s cart", act, req.MovieID, direction, req.Username)
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"msg": msg},
		)
		return
	}
	if !inserted {
		if response == nil {
			msg := fmt.Sprintf("%s is already in %s cart", req.MovieID, req.Username)
			if action == constants.DELETE {
				msg = fmt.Sprintf("%s was not in %s cart", req.MovieID, req.Username)
			}
			c.IndentedJSON(
				http.StatusNotModified,
				gin.H{"msg": msg},
			)
		} else {
			msg := fmt.Sprintf("Failed to add movie %s to %s cart", req.MovieID, req.Username)
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"msg": msg},
			)
		}

	} else {
		c.IndentedJSON(http.StatusAccepted, response)
	}
}

type UpdataeInventoryRequest struct {
	Username string   `json:"username"`
	MovieIDs []string `json:"movie_ids"`
}

func CheckoutEndpoint(c *gin.Context) {
	checkoutReturnHelper(c, memberRepo.Checkout)
}

func ReturnEndpoint(c *gin.Context) {
	checkoutReturnHelper(c, memberRepo.Return)
}

func checkoutReturnHelper(c *gin.Context, f func(context.Context, string, []string) ([]string, int, error)) {
	uir := UpdataeInventoryRequest{}
	err := c.BindJSON(&uir)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"msg": "Failed to unmarshall data into request"},
		)
		return
	}

	messages, moviesProcessed, err := f(c, uir.Username, uir.MovieIDs)
	if err != nil {
		msg := fmt.Sprintf("Failed to checkout %s\n", uir.Username)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": msg})
		return
	}

	status := http.StatusAccepted
	if moviesProcessed == 0 {
		status = http.StatusNotModified
	}
	c.IndentedJSON(
		status,
		gin.H{
			"messages":         messages,
			"movies_processed": moviesProcessed,
		},
	)
}

func GetCheckedOutMovies(c *gin.Context) {
	user, err := memberRepo.GetMemberByUsername(c, c.Param("username"), constants.CART)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve user cart"})
		return
	}
	movies, err := memberRepo.MovieRepo.GetMoviesByID(c, user.Checkedout, constants.CART)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"msg": fmt.Sprintf("Failed to retrieve checkedout movies from cloud for user %s", user.Username)})
		return
	}
	c.IndentedJSON(http.StatusOK, movies)
}

func SetMemberAPIChoiceEndpoint(c *gin.Context) {
	username := c.Param(constants.USERNAME)
	apiChoice := c.Query(constants.API_CHOICE)
	if username == "" || apiChoice == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Username and API choice are required"})
		return
	}
	if !(apiChoice == constants.REST_API || apiChoice == constants.GRAPHQL_API) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "Invalid API choice"})
		return
	}

	err := memberRepo.SetMemberAPIChoice(c, username, apiChoice)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("Failed to set API choice for user %s", username)})
		return
	}

	c.IndentedJSON(http.StatusAccepted, gin.H{"msg": fmt.Sprintf("API choice for user %s set to %s", username, apiChoice)})
}
