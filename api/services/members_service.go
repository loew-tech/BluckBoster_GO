package services

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/repos"
	"blockbuster/api/utils"
)

type MembersService struct {
	repo repos.MemberRepoInterface
}

var (
	instantiateServiceOnce sync.Once
	membersService         *MembersService
)

func GetMemberService() *MembersService {
	instantiateServiceOnce.Do(func() {
		membersService = &MembersService{repo: repos.NewMemberRepoWithDynamo()}
	})
	return membersService
}

func NewMemberServiceWithRepo(repo repos.MemberRepoInterface) *MembersService {
	return &MembersService{repo: repo}
}

func (s *MembersService) GetMember(c context.Context, username string, forCart bool) (data.Member, error) {
	member, err := s.repo.GetMemberByUsername(c, username, forCart)
	if err != nil {
		utils.LogError(fmt.Sprintf("failed to retrieve user %s", username), err)
		return data.Member{}, fmt.Errorf("failed to retrieve user %s", username)
	}
	return member, nil
}

func (s *MembersService) Login(c context.Context, username string) (data.Member, error) {
	member, err := s.repo.GetMemberByUsername(c, username, constants.NOT_CART)
	if err != nil || member.Username == "" {
		utils.LogError(fmt.Sprintf("user %s not found", username), err)
		return data.Member{}, fmt.Errorf("failed to login with user %s", username)
	}
	return member, nil
}

func (s *MembersService) GetCartIDs(c context.Context, username string) ([]string, error) {
	user, err := s.repo.GetMemberByUsername(c, username, constants.CART)
	if err != nil {
		utils.LogError(fmt.Sprintf("sser %s not found", username), err)
		return nil, fmt.Errorf("user %s not found", username)
	}
	return user.Cart, nil
}

func (s *MembersService) GetCartMovies(c context.Context, username string) ([]data.Movie, error) {
	movies, err := s.repo.GetCartMovies(c, username)
	if err != nil {
		utils.LogError(fmt.Sprintf("failed to retrieve cart movies for %s", username), err)
		return nil, fmt.Errorf("failed to retrieve cart movies for %s", username)
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
		utils.LogError("err updating cart", err)
		return false, errors.New("err updating cart")
	}
	return modified, nil
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
		utils.LogError(fmt.Sprintf("failed to update inventory for %s", username), err)
		return nil, 0, fmt.Errorf("failed to update inventory for %s", username)
	}
	return msgs, modifiedCount, nil
}

func (s *MembersService) GetCheckedOutMovies(c context.Context, username string) ([]data.Movie, error) {
	movies, err := s.repo.GetCheckedOutMovies(c, username)
	if err != nil {
		utils.LogError(fmt.Sprintf("failed to retrieve checked out movies for %s", username), err)
		return nil, fmt.Errorf("failed to retrieve checked out movies for %s", username)
	}
	return movies, nil
}

func (s *MembersService) SetAPIChoice(c context.Context, username, apiChoice string) error {
	if err := s.repo.SetMemberAPIChoice(c, username, apiChoice); err != nil {
		utils.LogError(fmt.Sprintf("failed to set %s apiChoice to %s", username, apiChoice), err)
		return fmt.Errorf("failed to set %s api selection to %s", username, apiChoice)
	}
	return nil
}

func (s *MembersService) UpdateMood(c context.Context, currentMood data.MovieMetrics, iteration int, movieIDs []string) (data.MovieMetrics, error) {
	mood, err := s.repo.UpdateMood(c, currentMood, iteration, movieIDs)
	if err != nil {
		return mood, utils.LogError("errs occurred while updatings mood", nil)
	}
	return mood, nil
}

func (s *MembersService) GetIniitialVotingSlate(c context.Context) ([]string, error) {
	movieIDs, err := s.repo.GetIniitialVotingSlate(c)
	if err != nil {
		return movieIDs, utils.LogError("errs occurred in getting initial voting slate", nil)
	}
	return movieIDs, nil
}

func (s *MembersService) IterateRecommendationVoting(c context.Context, currentMood data.MovieMetrics, iteration int, movieIDs []string) (data.MovieMetrics, []string, error) {
	mood, newMovieIDs, err := s.repo.IterateRecommendationVoting(c, currentMood, iteration, movieIDs)
	if err != nil {
		return mood, newMovieIDs, utils.LogError("iterated recommendation voting with errs", nil)
	}
	return mood, newMovieIDs, nil
}
