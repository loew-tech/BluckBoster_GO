package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	"blockbuster/api/repos"
)

type MembersHandler struct {
	repo repos.MemberRepoInterface
}

func NewMembersHandler() *MembersHandler {
	return &MembersHandler{
		repo: repos.NewMemberRepoWithDynamo(),
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
	username := c.Param(constants.USERNAME)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing username parameter"})
		return
	}
	member, err := h.repo.GetMemberByUsername(c, username, constants.NOT_CART)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve user"})
		return
	}
	if member.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("User %s not found", username)})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (h *MembersHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body"})
		return
	}
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing username in request body"})
		return
	}
	member, err := h.repo.GetMemberByUsername(c, req.Username, constants.NOT_CART)
	if err != nil || member.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("User %s not found", req.Username)})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (h *MembersHandler) GetCartIDs(c *gin.Context) {
	username := c.Param(constants.USERNAME)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing username parameter"})
		return
	}
	user, err := h.repo.GetMemberByUsername(c, username, constants.CART)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Cart not found"})
		return
	}
	c.JSON(http.StatusOK, user.Cart)
}

func (h *MembersHandler) GetCartMovies(c *gin.Context) {
	username := c.Param(constants.USERNAME)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing username parameter"})
		return
	}
	movies, err := h.repo.GetCartMovies(c, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve cart movies"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MembersHandler) AddToCart(c *gin.Context) {
	h.modifyCart(c, constants.ADD, false)
}

func (h *MembersHandler) RemoveFromCart(c *gin.Context) {
	h.modifyCart(c, constants.DELETE, false)
}

func (h *MembersHandler) modifyCart(c *gin.Context, action string, checkingOut bool) {
	var req struct {
		Username string `json:"username"`
		MovieID  string `json:"movie_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid cart request"})
		return
	}
	if req.MovieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing movie_id parameter"})
		return
	}

	inserted, response, err := h.repo.ModifyCart(c, req.Username, req.MovieID, action, checkingOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Cart modification failed"})
		return
	}
	if !inserted {
		c.JSON(http.StatusNotModified, gin.H{"msg": "No changes made to cart"})
		return
	}
	c.JSON(http.StatusAccepted, response)
}

func (h *MembersHandler) Checkout(c *gin.Context) {
	h.handleInventoryAction(c, h.repo.Checkout)
}

func (h *MembersHandler) Return(c *gin.Context) {
	h.handleInventoryAction(c, h.repo.Return)
}

func (h *MembersHandler) handleInventoryAction(
	c *gin.Context,
	f func(context.Context, string, []string) ([]string, int, error),
) {
	var req struct {
		Username string   `json:"username"`
		MovieIDs []string `json:"movie_ids"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body"})
		return
	}
	if len(req.MovieIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No movie_ids provided"})
		return
	}
	msgs, count, err := f(c, req.Username, req.MovieIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to process inventory"})
		return
	}
	status := http.StatusAccepted
	if count == 0 {
		status = http.StatusNotModified
	}
	c.JSON(status, gin.H{"messages": msgs, "movies_processed": count})
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
