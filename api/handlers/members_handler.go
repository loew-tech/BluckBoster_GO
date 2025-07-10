package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	"blockbuster/api/repos"
	"blockbuster/api/services"
)

type MembersHandler struct {
	repo    repos.MemberRepoInterface
	service *services.MembersService
}

func NewMembersHandler() *MembersHandler {
	return &MembersHandler{
		repo:    repos.NewMemberRepoWithDynamo(),
		service: services.GetMemberService(),
	}
}

func (h *MembersHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/members/:username", h.GetMember)
	rg.POST("/members/login", h.Login)
	rg.GET("/members/:username/cart", h.GetCartMovies)
	rg.GET("/members/cart/:username", h.GetCartMovies)
	rg.PUT("/members/cart", h.AddToCart)
	rg.PUT("/members/cart/remove", h.RemoveFromCart)
	rg.GET("/members/:username/cart/ids", h.GetCartIDs)
	rg.POST("/members/checkout", h.Checkout)
	rg.POST("/members/return", h.Return)
	rg.GET("/members/:username/checkedout", h.GetCheckedOutMovies)
	rg.PUT("/members/:username", h.SetAPIChoice)
}

func (h *MembersHandler) GetMember(c *gin.Context) {
	status, member, err := h.service.GetMember(c)
	if err != nil {
		c.JSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, member)
}

func (h *MembersHandler) Login(c *gin.Context) {
	status, member, err := h.service.Login(c)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(status, member)
}

func (h *MembersHandler) GetCartIDs(c *gin.Context) {
	status, cartIDs, err := h.service.GetCartIDs(c)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(status, cartIDs)
}

func (h *MembersHandler) GetCartMovies(c *gin.Context) {
	status, movies, err := h.service.GetCartMovies(c)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, movies)
}

func (h *MembersHandler) AddToCart(c *gin.Context) {
	status, err := h.service.AddToCart(c)
	if err != nil {
		c.JSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, gin.H{"msg": "success"})
}

func (h *MembersHandler) RemoveFromCart(c *gin.Context) {
	status, err := h.service.RemoveFromCart(c)
	if err != nil {
		c.JSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, gin.H{"msg": "success"})
}

func (h *MembersHandler) Checkout(c *gin.Context) {
	status, msgs, modifiedCount, err := h.service.Checkout(c)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(status, gin.H{"msgs": msgs, "Modified": modifiedCount})
}

func (h *MembersHandler) Return(c *gin.Context) {
	status, msgs, modifiedCount, err := h.service.Return(c)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(status, gin.H{"msgs": msgs, "Modified": modifiedCount})
}

func (h *MembersHandler) GetCheckedOutMovies(c *gin.Context) {
	username := c.Param(constants.USERNAME)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing username parameter"})
		return
	}
	_, err := h.repo.GetMemberByUsername(c, username, constants.CART)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve user"})
		return
	}
	movies, err := h.repo.GetCheckedoutMovies(c, username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"msg": "Failed to retrieve checked out movies"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MembersHandler) SetAPIChoice(c *gin.Context) {
	username := c.Param(constants.USERNAME)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing username parameter"})
		return
	}
	apiChoice := c.Query(constants.API_CHOICE)
	if apiChoice == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing api_choice parameter"})
		return
	}
	if apiChoice != constants.REST_API && apiChoice != constants.GRAPHQL_API {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid API choice"})
		return
	}
	if err := h.repo.SetMemberAPIChoice(c, username, apiChoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update API choice"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"msg": fmt.Sprintf("API choice set to %s", apiChoice)})
}
