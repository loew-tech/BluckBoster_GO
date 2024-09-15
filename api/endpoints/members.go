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

type Username struct {
	Username string `json:"username"`
}

func MemberLoginEndpoint(c *gin.Context) {
	lr := Username{}
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
	movies, err := memberRepo.GetCartIDs(c.Param("username"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve user cart"})
	} else {
		c.IndentedJSON(http.StatusAccepted, movies)
	}
}

func GetCartMoviesEndpoint(c *gin.Context) {
	movies, err := memberRepo.GetCartMovies(c.Param("username"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve cart ids"})
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
	inserted, response, err := memberRepo.ModifyCart(req.Username, req.MovieID, db.ADD)
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
	removed, response, err := memberRepo.ModifyCart(req.Username, req.MovieID, db.DELETE)
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
