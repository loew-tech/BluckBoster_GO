package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/repos"
	"blockbuster/api/utils"
)

// @TODO: make err handling consistent with movies_service
type MembersService struct {
	repo repos.MemberRepoInterface
}

var (
	instantiateServiceOnce sync.Once
	service                *MembersService
)

func GetMemberService() *MembersService {
	instantiateServiceOnce.Do(func() {
		service = &MembersService{repo: repos.NewMemberRepoWithDynamo()}
	})
	return service
}

func (s *MembersService) GetMember(c *gin.Context) (int, data.Member, error) {
	username := c.Param(constants.USERNAME)
	if username == "" {
		return http.StatusBadRequest, data.Member{}, errors.New("")
	}
	member, err := s.repo.GetMemberByUsername(c, username, constants.NOT_CART)
	if err != nil {
		return http.StatusInternalServerError, data.Member{}, utils.LogError(fmt.Sprintf("Failed to retrieve user %s", username), err)
	}
	return http.StatusOK, member, nil
}

func (s *MembersService) Login(c *gin.Context) (int, data.Member, error) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" {
		return http.StatusBadRequest, data.Member{}, utils.LogError("Invalid login request body", nil)
	}
	member, err := s.repo.GetMemberByUsername(c, req.Username, constants.NOT_CART)
	if err != nil || member.Username == "" {
		return http.StatusNotFound, data.Member{}, utils.LogError(fmt.Sprintf("User %s not found", req.Username), nil)
	}
	return http.StatusOK, member, nil
}

func (s *MembersService) GetCartIDs(c *gin.Context) (int, []string, error) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	user, err := s.repo.GetMemberByUsername(c, username, constants.CART)
	if err != nil {
		return http.StatusNotFound, nil, utils.LogError(fmt.Sprintf("User %s not found", username), err)
	}
	return http.StatusOK, user.Cart, nil
}

func (s *MembersService) GetCartMovies(c *gin.Context) (int, []data.Movie, error) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	movies, err := s.repo.GetCartMovies(c, username)
	if err != nil {
		return http.StatusInternalServerError, nil, utils.LogError(fmt.Sprintf("Failed to retrieve cart movies for %s", username), err)
	}
	return http.StatusOK, movies, nil
}

func (s *MembersService) AddToCart(c *gin.Context) (int, error) {
	return s.modifyCart(c, constants.ADD, constants.NOT_CHECKOUT)
}

func (s *MembersService) RemoveFromCart(c *gin.Context) (int, error) {
	return s.modifyCart(c, constants.DELETE, constants.NOT_CHECKOUT)
}

func (s *MembersService) modifyCart(c *gin.Context, action string, checkingOut bool) (int, error) {
	var req struct {
		Username string `json:"username"`
		MovieID  string `json:"movie_id"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || req.MovieID == "" {
		return http.StatusBadRequest, utils.LogError("Invalid modify cart request", nil)
	}
	modified, _, err := s.repo.ModifyCart(c, req.Username, req.MovieID, action, checkingOut)
	if err != nil {
		return http.StatusInternalServerError, utils.LogError("err updating cart", err)
	}
	status := http.StatusOK
	if modified {
		status = http.StatusAccepted
	}
	return status, err
}

func (s MembersService) Checkout(c *gin.Context) (int, []string, int, error) {
	return s.handleInventoryAction(c, s.repo.Checkout)
}

func (s *MembersService) Return(c *gin.Context) (int, []string, int, error) {
	return s.handleInventoryAction(c, s.repo.Return)
}

func (s *MembersService) handleInventoryAction(
	c *gin.Context,
	f func(context.Context, string, []string) ([]string, int, error),

) (int, []string, int, error) {
	var req struct {
		Username string   `json:"username"`
		MovieIDs []string `json:"movie_ids"`
	}
	if err := c.BindJSON(&req); err != nil || req.Username == "" || len(req.MovieIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body"})
		return http.StatusBadRequest, nil, 0, utils.LogError("Invalid inventory request", err)
	}
	msgs, modifiedCount, err := f(c, req.Username, req.MovieIDs)
	if err != nil {
		return http.StatusInternalServerError, nil, 0, fmt.Errorf("failed to update inventory for %s", req.Username)
	}
	status := http.StatusOK
	if 0 < modifiedCount {
		status = http.StatusAccepted
	}
	return status, msgs, modifiedCount, nil
}

func (s *MembersService) GetCheckedOutMovies(c *gin.Context) (int, []data.Movie, error) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	movies, err := s.repo.GetCheckedOutMovies(c, username)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to retrieve checked out movies for %s", username)
	}
	return http.StatusAccepted, movies, nil
}

func (s *MembersService) SetAPIChoice(c *gin.Context) (int, string, error) {
	username, err := utils.GetStringArg(c.Params, constants.USERNAME)
	if err != nil {
		return http.StatusBadRequest, "", errors.New("missing param 'username'")
	}
	apiChoice := c.Query(constants.API_CHOICE)
	if apiChoice == "" || (apiChoice != constants.REST_API && apiChoice != constants.GRAPHQL_API) {
		return http.StatusBadRequest, "", fmt.Errorf("invalid api choice for %s. Selected '%s' but must be '%s' or '%s'", username, apiChoice, constants.REST_API, constants.GRAPHQL_API)
	}
	if err := s.repo.SetMemberAPIChoice(c, username, apiChoice); err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("failed to set %s api selection to %s", username, apiChoice)
	}
	return http.StatusAccepted, apiChoice, nil
}
