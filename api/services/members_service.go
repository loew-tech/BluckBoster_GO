package services

import (
	"context"
	"fmt"
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

func (s *MembersService) GetMember(c context.Context, username string) (data.Member, error) {
	member, err := s.repo.GetMemberByUsername(c, username, constants.NOT_CART)
	if err != nil {
		return data.Member{}, utils.LogError(fmt.Sprintf("Failed to retrieve user %s", username), err)
	}
	return member, nil
}

func (s *MembersService) Login(c context.Context, username string) (data.Member, error) {
	member, err := s.repo.GetMemberByUsername(c, username, constants.NOT_CART)
	if err != nil || member.Username == "" {
		return data.Member{}, utils.LogError(fmt.Sprintf("User %s not found", username), nil)
	}
	return member, nil
}

func (s *MembersService) GetCartIDs(c context.Context, username string) ([]string, error) {
	user, err := s.repo.GetMemberByUsername(c, username, constants.CART)
	if err != nil {
		return nil, utils.LogError(fmt.Sprintf("User %s not found", username), err)
	}
	return user.Cart, nil
}

func (s *MembersService) GetCartMovies(c context.Context, username string) ([]data.Movie, error) {
	movies, err := s.repo.GetCartMovies(c, username)
	if err != nil {
		return nil, utils.LogError(fmt.Sprintf("Failed to retrieve cart movies for %s", username), err)
	}
	return movies, nil
}

func (s *MembersService) AddToCart(c context.Context, username, movieID string) (bool, error) {
	return s.modifyCart(c, username, movieID, constants.ADD, constants.NOT_CHECKOUT)
}

func (s *MembersService) RemoveFromCart(c context.Context, username, movieID string) (bool, error) {
	return s.modifyCart(c, username, movieID, constants.DELETE, constants.NOT_CHECKOUT)
}

func (s *MembersService) modifyCart(c context.Context, username, movieID, action string, checkingOut bool) (bool, error) {
	modified, err := s.repo.ModifyCart(c, username, movieID, action, checkingOut)
	if err != nil {
		return false, utils.LogError("err updating cart", err)
	}
	return modified, err
}

func (s MembersService) Checkout(c context.Context, username string, movieIDs []string) ([]string, int, error) {
	return s.handleInventoryAction(c, username, movieIDs, s.repo.Checkout)
}

func (s *MembersService) Return(c context.Context, username string, movieIDs []string) ([]string, int, error) {
	return s.handleInventoryAction(c, username, movieIDs, s.repo.Return)
}

func (s *MembersService) handleInventoryAction(
	c context.Context, username string, movieIDs []string,
	f func(context.Context, string, []string) ([]string, int, error),
) ([]string, int, error) {
	msgs, modifiedCount, err := f(c, username, movieIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to update inventory for %s", username)
	}
	return msgs, modifiedCount, nil
}

func (s *MembersService) GetCheckedOutMovies(c context.Context, username string) ([]data.Movie, error) {
	movies, err := s.repo.GetCheckedOutMovies(c, username)
	if err != nil {
		return nil, utils.LogError(fmt.Sprintf("failed to retrieve checked out movies for %s", username), err)
	}
	return movies, nil
}

func (s *MembersService) SetAPIChoice(c context.Context, username, apiChoice string) error {
	if err := s.repo.SetMemberAPIChoice(c, username, apiChoice); err != nil {
		return fmt.Errorf("failed to set %s api selection to %s", username, apiChoice)
	}
	return nil
}
