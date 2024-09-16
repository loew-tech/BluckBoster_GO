package endpoints

import (
	"blockbuster/api/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var memberRepo = db.NewMembersRepo()

func GetMemberEndpoint(c *gin.Context) {
	found, member, err := memberRepo.GetMemberByUsername(c.Param("username"), db.NOT_CART)
	if err != nil {
		// @TODO sometimes should this be a 502?
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
		c.IndentedJSON(http.StatusAccepted, user.Cart)
	}
}

func GetCartMoviesEndpoint(c *gin.Context) {
	movies, err := memberRepo.GetCartMovies(c.Param("username"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve cart movies"})
	} else {
		c.IndentedJSON(http.StatusAccepted, movies)
	}
}

type ModifyCartRequest struct {
	Username string `json:"username"`
	MovieID  string `json:"movie_id"`
}

func AddToCartEndpoint(c *gin.Context) {
	req := ModifyCartRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"msg": "Bad Request for ModifyCart"},
		)
		return
	}
	inserted, response, err := memberRepo.ModifyCart(req.Username, req.MovieID, db.ADD, db.NOT_CHECKOUT)
	if err != nil {
		msg := fmt.Sprintf("Error adding movie %s to %s cart", req.MovieID, req.Username)
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"msg": msg},
		)
		return
	}
	if !inserted {
		if response == nil {
			msg := fmt.Sprintf("Movie %s already in %s cart", req.MovieID, req.Username)
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

func RemoveFromCartEndpoint(c *gin.Context) {
	req := ModifyCartRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"msg": "Bad Request for ModifyCart"},
		)
		return
	}
	removed, response, err := memberRepo.ModifyCart(req.Username, req.MovieID, db.DELETE, db.NOT_CHECKOUT)
	if err != nil {
		msg := fmt.Sprintf("Error removing %s from %s cart", req.MovieID, req.Username)
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"msg": msg},
		)
		return
	}
	if !removed {
		if response == nil {
			msg := fmt.Sprintf("%s was not in %s cart", req.MovieID, req.Username)
			c.IndentedJSON(
				http.StatusNotModified,
				gin.H{"msg": msg},
			)
		} else {
			msg := fmt.Sprintf("Failed to remove movie %s from %s cart", req.MovieID, req.Username)
			c.IndentedJSON(
				http.StatusInternalServerError,
				gin.H{"msg": msg},
			)
		}

	} else {
		c.IndentedJSON(http.StatusAccepted, response)
	}
}

func CheckoutEndpoint(c *gin.Context) {
	fmt.Println("In checkout endpoint")
	un := UsernameReq{}
	err := c.BindJSON(&un)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"msg": "Failed to unmarshall data into request"},
		)
		return
	}

	messages, rented, err := memberRepo.Checkout(un.Username)
	errMsg := fmt.Sprintf("Failed to checkout %s\n%s", un, err)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, errMsg)
	} else if rented == 0 {
		c.IndentedJSON(http.StatusNotModified, errMsg)
	} else {
		c.IndentedJSON(
			http.StatusAccepted,
			gin.H{
				"messages": messages,
				"rented":   rented,
			},
		)
	}
}
