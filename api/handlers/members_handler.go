package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	"blockbuster/api/services"
	"blockbuster/api/utils"
)

type MembersHandler struct {
	service services.MembersServiceInterface
}

func NewMembersHandler() *MembersHandler {
	return &MembersHandler{
		service: services.GetMemberService(),
	}
}

func NewMembersHandlerWithService(service services.MembersServiceInterface) *MembersHandler {
	return &MembersHandler{
		service: service,
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
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	member, err := h.service.GetMember(c.Request.Context(), username, constants.NOT_CART)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("error occured retrieving user %s", username)})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (h *MembersHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid value for field 'username'"})
		return
	}
	log.Printf("attempting login with username %s\n", req.Username)
	member, err := h.service.Login(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("error logging in as user %s", req.Username)})
		return
	}
	if member.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("failed logging in as user %s", req.Username)})
		return
	}
	c.IndentedJSON(http.StatusOK, member)
}

func (h *MembersHandler) GetCartIDs(c *gin.Context) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	cartIDs, err := h.service.GetCartIDs(c.Request.Context(), username)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, cartIDs)
}

func (h *MembersHandler) GetCartMovies(c *gin.Context) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	movies, err := h.service.GetCartMovies(c.Request.Context(), username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("Failed to retrieve cart movies for %s", username)})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MembersHandler) AddToCart(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		MovieID  string `json:"movie_id"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || req.MovieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid modify cart request; username and movie_id must be provided"})
		return
	}
	modified, err := h.service.AddToCart(c.Request.Context(), req.Username, req.MovieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if modified {
		c.JSON(http.StatusAccepted, gin.H{"msg": "success"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (h *MembersHandler) RemoveFromCart(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		MovieID  string `json:"movie_id"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || req.MovieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body. 'username' and 'movie_id' must be provided"})
		return
	}
	modified, err := h.service.RemoveFromCart(c.Request.Context(), req.Username, req.MovieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if modified {
		c.JSON(http.StatusAccepted, gin.H{"msg": "success"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (h *MembersHandler) Checkout(c *gin.Context) {
	var req struct {
		Username string   `json:"username"`
		MovieIDs []string `json:"movie_ids"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || len(req.MovieIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body. Requires 'username' and 'movie_ids'"})
		return
	}
	msgs, modifiedCount, err := h.service.Checkout(c.Request.Context(), req.Username, req.MovieIDs)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{"msgs": msgs, "Modified": modifiedCount})
}

func (h *MembersHandler) Return(c *gin.Context) {
	var req struct {
		Username string   `json:"username"`
		MovieIDs []string `json:"movie_ids"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || len(req.MovieIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body. Requires 'username' and 'movie_ids'"})
		return
	}
	msgs, modifiedCount, err := h.service.Return(c.Request.Context(), req.Username, req.MovieIDs)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{"msgs": msgs, "Modified": modifiedCount})
}

func (h *MembersHandler) GetCheckedOutMovies(c *gin.Context) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	movies, err := h.service.GetCheckedOutMovies(c.Request.Context(), username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MembersHandler) SetAPIChoice(c *gin.Context) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	apiChoice := c.Query(constants.API_CHOICE)
	if apiChoice == "" || (apiChoice != constants.REST_API && apiChoice != constants.GRAPHQL_API) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid api selection '%s'; must be 'REST' or 'GraphQL'", apiChoice)})
		return
	}
	if err = h.service.SetAPIChoice(c.Request.Context(), username, apiChoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("API choice set to %s", apiChoice)})
}
