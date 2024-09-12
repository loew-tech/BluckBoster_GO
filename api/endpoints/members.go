package endpoints

import (
	"blockbuster/api/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var memberRepo = db.NewMembersRepo()

func GetMemberEndpoint(c *gin.Context) {
	found, member, err := memberRepo.GetMemberByUsername(c.Param("username"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve user"})
	} else {
		if found {
			c.IndentedJSON(http.StatusOK, member)
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to find user"})
		}
	}
}

type LoginRequest struct {
	Username string `json:"username"`
}

func MemberLoginEndpoint(c *gin.Context) {
	lr := LoginRequest{}
	err := c.BindJSON(&lr)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"msg": "Failed to unmarshall data into request"},
		)
		return
	}
	found, member, err := memberRepo.GetMemberByUsername(lr.Username)
	if err != nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"msg": "Failed to retrieve user"},
		)
	} else {
		if found {
			c.IndentedJSON(http.StatusOK, member)
		} else {
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"msg": "Failed to find user"},
			)
		}
	}
}

func GetCartIDsEndpoint(c *gin.Context) {
	movies, err := memberRepo.GetCartIDs(c.Param("username"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve user cart"})
	} else {
		for m := range movies {
			fmt.Println(m)
		}
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
			gin.H{"msg": "Bad Request for AddToCart"},
		)
		return
	}
	inserted, response, err := memberRepo.AddToCart(req.Username, req.MovieID)
	if err != nil {
		msg := fmt.Sprintf("Error adding movie %s to %s cart", req.MovieID, req.Username)
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"msg": msg},
		)
		return
	}
	if !inserted {
		if response == nil {
			msg := fmt.Sprintf("Movie %s already in %s cart", req.MovieID, req.Username)
			c.IndentedJSON(
				http.StatusNotFound,
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
			gin.H{"msg": "Bad Request for AddToCart"},
		)
		return
	}
	removed, response, err := memberRepo.RemoveFromCart(req.Username, req.MovieID)
	if err != nil {
		msg := fmt.Sprintf("Error adding movie %s to %s cart", req.MovieID, req.Username)
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"msg": msg},
		)
		return
	}
	if !removed {
		if response == nil {
			msg := fmt.Sprintf("Movie %s was not %s cart", req.MovieID, req.Username)
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"msg": msg},
			)
		} else {
			msg := fmt.Sprintf("Failed to remove movie %s from %s cart", req.MovieID, req.Username)
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"msg": msg},
			)
		}

	} else {
		c.IndentedJSON(http.StatusAccepted, response)
	}
}
