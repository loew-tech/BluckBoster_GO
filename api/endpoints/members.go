package endpoints

import (
	"blockbuster/api/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

var memberRepo = db.NewMembersRepo()

func GetMemberEndpoint(c *gin.Context) {
	member, err := memberRepo.GetMemberByUsername(c.Param("username"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve user"})
	} else {
		c.IndentedJSON(http.StatusOK, member)
	}
}