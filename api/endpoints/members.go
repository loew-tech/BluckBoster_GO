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

type AddToCartRequest struct {
	Username string `json:"username"`
	LastName string `json:"last_name"`
	MovieID  string `json:"movie_id"`
}

func AddToCartEndpoint(c *gin.Context) {
	req := AddToCartRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"msg": "Bad Request for AddToCart"},
		)
		return
	}
	inserted, response, err := memberRepo.AddToCart(req.Username, req.MovieID, req.LastName)
	if err != nil {
		msg := fmt.Sprintf("Error adding movie %s to %s cart", req.MovieID, req.Username)
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"msg": msg},
		)
		return
	}
	if !inserted {
		msg := fmt.Sprintf("Failed to add movie %s to %s cart", req.MovieID, req.Username)
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"msg": msg},
		)
	} else {
		c.IndentedJSON(http.StatusAccepted, response)
	}
}
