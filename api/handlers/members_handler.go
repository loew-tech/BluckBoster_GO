package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	"blockbuster/api/services"
	"blockbuster/api/utils"
)

type MembersHandler struct {
	service *services.MembersService
}

func NewMembersHandler() *MembersHandler {
	return &MembersHandler{
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
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	status, member, err := h.service.GetMember(c.Request.Context(), username)
	if err != nil {
		c.JSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, member)
}

func (h *MembersHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid value for field 'username'"})
		return
	}
	status, member, err := h.service.Login(c.Request.Context(), req.Username)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(status, member)
}

func (h *MembersHandler) GetCartIDs(c *gin.Context) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	status, cartIDs, err := h.service.GetCartIDs(c.Request.Context(), username)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(status, cartIDs)
}

func (h *MembersHandler) GetCartMovies(c *gin.Context) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		// return http.StatusBadRequest, nil, err
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	status, movies, err := h.service.GetCartMovies(c.Request.Context(), username)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, movies)
}

func (h *MembersHandler) AddToCart(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		MovieID  string `json:"movie_id"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || req.MovieID == "" {
		// return http.StatusBadRequest, utils.LogError("Invalid modify cart request", nil)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid modify cart request; username and movie_id must be provided"})
		return
	}
	status, err := h.service.AddToCart(c, req.Username, req.MovieID)
	if err != nil {
		c.JSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, gin.H{"msg": "success"})
}

func (h *MembersHandler) RemoveFromCart(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		MovieID  string `json:"movie_id"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || req.MovieID == "" {
		// return http.StatusBadRequest, utils.LogError("Invalid modify cart request", nil)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid modify cart request; username and movie_id must be provided"})
		return
	}
	status, err := h.service.RemoveFromCart(c.Request.Context(), req.Username, req.MovieID)
	if err != nil {
		c.JSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, gin.H{"msg": "success"})
}

func (h *MembersHandler) Checkout(c *gin.Context) {
	var req struct {
		Username string   `json:"username"`
		MovieIDs []string `json:"movie_ids"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || len(req.MovieIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body. Requires username and movie_ids"})
		return
	}
	status, msgs, modifiedCount, err := h.service.Checkout(c, req.Username, req.MovieIDs)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(status, gin.H{"msgs": msgs, "Modified": modifiedCount})
}

func (h *MembersHandler) Return(c *gin.Context) {
	var req struct {
		Username string   `json:"username"`
		MovieIDs []string `json:"movie_ids"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || len(req.MovieIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body. Requires username and movie_ids"})
		return
	}
	status, msgs, modifiedCount, err := h.service.Return(c, req.Username, req.MovieIDs)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.IndentedJSON(status, gin.H{"msgs": msgs, "Modified": modifiedCount})
}

func (h *MembersHandler) GetCheckedOutMovies(c *gin.Context) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
		// return http.StatusBadRequest, nil, err
	}
	status, movies, err := h.service.GetCheckedOutMovies(c.Request.Context(), username)
	if err != nil {
		c.IndentedJSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, movies)
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
	status, apiChoice, err := h.service.SetAPIChoice(c.Request.Context(), username, apiChoice)
	if err != nil {
		c.JSON(status, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(status, gin.H{"msg": fmt.Sprintf("API choice set to %s", apiChoice)})
}
