package endpoints

import (
	"blockbuster/api/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var memberRepo = db.NewMembersRepo(GetDynamoClient())

func GetMemberEndpoint(c *gin.Context) {
	found, member, err := memberRepo.GetMemberByUsername(c.Param("username"), db.NOT_CART)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve user"})
	} else {
		if found {
			c.IndentedJSON(http.StatusOK, member)
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to find user"})
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
	found, member, err := memberRepo.GetMemberByUsername(un.Username, db.NOT_CART)
	if err != nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"msg": "Error retrieving user"},
		)
		return
	}
	if found {
		c.IndentedJSON(http.StatusOK, member)
	} else {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"msg": "Failed to find user"},
		)
	}
}

func GetCartIDsEndpoint(c *gin.Context) {
	_, user, err := memberRepo.GetMemberByUsername(c.Param("username"), db.CART)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve user cart"})
	} else {
		c.IndentedJSON(http.StatusOK, user.Cart)
	}
}

func GetCartMoviesEndpoint(c *gin.Context) {
	movies, err := memberRepo.GetCartMovies(c.Param("username"))
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
	cartHelper(c, db.ADD, false)
}

func RemoveFromCartEndpoint(c *gin.Context) {
	cartHelper(c, db.DELETE, false)
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

	inserted, response, err := memberRepo.ModifyCart(req.Username, req.MovieID, action, checkingOut)
	if err != nil {
		act, direction := "adding", "to"
		if action == db.DELETE {
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
			if action == db.DELETE {
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

func checkoutReturnHelper(c *gin.Context, f func(string, []string) ([]string, int, error)) {
	uir := UpdataeInventoryRequest{}
	err := c.BindJSON(&uir)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"msg": "Failed to unmarshall data into request"},
		)
		return
	}

	messages, moviesProcessed, err := f(uir.Username, uir.MovieIDs)
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
	_, user, err := memberRepo.GetMemberByUsername(c.Param("username"), db.CART)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve user cart"})
	}
	_, movies, err := memberRepo.MovieRepo.GetMoviesByID(user.Checkedout, db.CART)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"msg": "Failed to retrieve movies from cloud"})
	} else {
		c.IndentedJSON(http.StatusOK, movies)
	}
}
