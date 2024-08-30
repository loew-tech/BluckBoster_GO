package endpoints

import (
	"blockbuster/api/db"
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
    username string
}

func MemberLoginEndpoint(c *gin.Context) {
	lr := LoginRequest{}
	c.BindJSON(&lr)
	found, member, err := memberRepo.GetMemberByUsername(lr.username)
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

type AddToCartRequest struct {
	username string
	lastName string
	movieID string
}

func AddToCartEndpoint(c *gin.Context) {
	cr := AddToCartRequest{}
	c.BindJSON(cr)
	// @TODO: finish endpopint
}