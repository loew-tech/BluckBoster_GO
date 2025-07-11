package services

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/repos"
	"blockbuster/api/utils"
)

// @TODO: write interfaces
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

func (s *MembersService) GetMember(c context.Context, username string) (int, data.Member, error) {
	member, err := s.repo.GetMemberByUsername(c, username, constants.NOT_CART)
	if err != nil {
		return http.StatusInternalServerError, data.Member{}, utils.LogError(fmt.Sprintf("Failed to retrieve user %s", username), err)
	}
	return http.StatusOK, member, nil
}

func (s *MembersService) Login(c context.Context, username string) (int, data.Member, error) {
	member, err := s.repo.GetMemberByUsername(c, username, constants.NOT_CART)
	// @TODO: could just return member and err and do this check in handler... or err won't be null
	if err != nil || member.Username == "" {
		return http.StatusNotFound, data.Member{}, utils.LogError(fmt.Sprintf("User %s not found", username), nil)
	}
	return http.StatusOK, member, nil
}

func (s *MembersService) GetCartIDs(c context.Context, username string) (int, []string, error) {
	user, err := s.repo.GetMemberByUsername(c, username, constants.CART)
	// @TODO: could just return nil and err and handle status in handler
	if err != nil {
		return http.StatusNotFound, nil, utils.LogError(fmt.Sprintf("User %s not found", username), err)
	}
	return http.StatusOK, user.Cart, nil
}

func (s *MembersService) GetCartMovies(c context.Context, username string) (int, []data.Movie, error) {
	movies, err := s.repo.GetCartMovies(c, username)
	// @TODO: return movies and err
	if err != nil {
		return http.StatusInternalServerError, nil, utils.LogError(fmt.Sprintf("Failed to retrieve cart movies for %s", username), err)
	}
	return http.StatusOK, movies, nil
}

func (s *MembersService) AddToCart(c context.Context, username, movieID string) (int, error) {
	return s.modifyCart(c, constants.ADD, username, movieID, constants.NOT_CHECKOUT)
}

func (s *MembersService) RemoveFromCart(c context.Context, username, movieID string) (int, error) {
	return s.modifyCart(c, constants.DELETE, username, movieID, constants.NOT_CHECKOUT)
}

func (s *MembersService) modifyCart(c context.Context, username, movieID, action string, checkingOut bool) (int, error) {
	// @TODO: return modifed, err
	modified, _, err := s.repo.ModifyCart(c, username, movieID, action, checkingOut)
	if err != nil {
		return http.StatusInternalServerError, utils.LogError("err updating cart", err)
	}
	status := http.StatusOK
	if modified {
		status = http.StatusAccepted
	}
	return status, err
}

func (s MembersService) Checkout(c context.Context, username string, movieIDs []string) (int, []string, int, error) {
	return s.handleInventoryAction(c, username, movieIDs, s.repo.Checkout)
}

func (s *MembersService) Return(c context.Context, username string, movieIDs []string) (int, []string, int, error) {
	return s.handleInventoryAction(c, username, movieIDs, s.repo.Return)
}

// @TODO: return bool (modifiedCount > 0) and err
func (s *MembersService) handleInventoryAction(
	c context.Context, username string, movieIDs []string,
	f func(context.Context, string, []string) ([]string, int, error),

) (int, []string, int, error) {
	msgs, modifiedCount, err := f(c, username, movieIDs)
	if err != nil {
		return http.StatusInternalServerError, nil, 0, fmt.Errorf("failed to update inventory for %s", username)
	}
	status := http.StatusOK
	if 0 < modifiedCount {
		status = http.StatusAccepted
	}
	return status, msgs, modifiedCount, nil
}

// @TODO: modify return
func (s *MembersService) GetCheckedOutMovies(c context.Context, username string) (int, []data.Movie, error) {
	movies, err := s.repo.GetCheckedOutMovies(c, username)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("Failed to retrieve checked out movies for %s", username)
	}
	return http.StatusAccepted, movies, nil
}

// @TODO: modify return
func (s *MembersService) SetAPIChoice(c context.Context, username, apiChoice string) (int, string, error) {
	if err := s.repo.SetMemberAPIChoice(c, username, apiChoice); err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("failed to set %s api selection to %s", username, apiChoice)
	}
	return http.StatusAccepted, apiChoice, nil
}
