package repos

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	api_cache "blockbuster/api/api_cache"
	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/utils"
)

const membersTableName = "BluckBoster_members"

type MemberRepo struct {
	client            DynamoClientInterface
	tableName         string
	movieRepo         ReadWriteMovieRepo
	centroids         api_cache.CentroidCacheInterface
	centroidsToMovies api_cache.CentroidsToMoviesCacheInterface
	movieMetricCache  map[string]data.MovieMetrics
	randGen           rand.Rand
}

func NewMembersRepo(client DynamoClientInterface, movieRepo ReadWriteMovieRepo, centroidsCache api_cache.CentroidCacheInterface, centroidsToMovies api_cache.CentroidsToMoviesCacheInterface) MemberRepoInterface {
	source := rand.NewSource(time.Now().UnixNano())
	return &MemberRepo{
		client:            client,
		tableName:         membersTableName,
		movieRepo:         movieRepo,
		centroids:         centroidsCache,
		centroidsToMovies: centroidsToMovies,
		movieMetricCache:  make(map[string]data.MovieMetrics),
		randGen:           *rand.New(source),
	}
}

func (r *MemberRepo) GetMemberByUsername(ctx context.Context, username string, cartOnly bool) (data.Member, error) {
	member := data.Member{}

	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			constants.USERNAME: &types.AttributeValueMemberS{Value: username},
		},
		TableName: &r.tableName,
	}

	if cartOnly {
		expr := "username, cart, checked_out, #t"
		input.ProjectionExpression = &expr
		input.ExpressionAttributeNames = map[string]string{"#t": constants.TYPE}
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return member, utils.LogError("fetching user from cloud", err)
	}
	if result.Item == nil {
		return member, utils.LogError(fmt.Sprintf("user %s not found", username), errors.New("item not found"))
	}

	err = attributevalue.UnmarshalMap(result.Item, &member)
	if err != nil {
		return member, utils.LogError("unmarshalling user data", err)
	}
	return member, nil
}

func (r *MemberRepo) GetCartMovies(ctx context.Context, username string) ([]data.Movie, error) {
	user, err := r.GetMemberByUsername(ctx, username, constants.CART)
	if err != nil {
		return nil, utils.LogError("fetching cart movie IDs", err)
	}
	if len(user.Cart) == 0 {
		return []data.Movie{}, nil
	}
	return r.movieRepo.GetMoviesByID(ctx, user.Cart, constants.CART)
}

func (r *MemberRepo) GetCheckedOutMovies(ctx context.Context, username string) ([]data.Movie, error) {
	if username == "" {
		return nil, errors.New("username is required to get checkout moves")
	}
	user, err := r.GetMemberByUsername(ctx, username, constants.CART)
	if err != nil {
		return nil, utils.LogError(fmt.Sprintf("Failed to get user for username %s", username), err)
	}
	movies, err := r.movieRepo.GetMoviesByID(ctx, user.Checkedout, constants.CART)
	if err != nil {
		return nil, utils.LogError(fmt.Sprintf("Error in fetching checkedout movies for %s", username), err)
	}
	return movies, nil
}

func (r *MemberRepo) ModifyCart(ctx context.Context, username, movieID, updateKey string, checkingOut bool) (bool, error) {
	name, err := attributevalue.Marshal(username)
	if err != nil {
		return false, utils.LogError("marshalling username", err)
	}

	expr, attrs := buildCartUpdateExpr(movieID, updateKey, checkingOut)
	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       map[string]types.AttributeValue{constants.USERNAME: name},
		ExpressionAttributeValues: attrs,
		ReturnValues:              types.ReturnValueUpdatedNew,
		UpdateExpression:          &expr,
	}
	return r.updateMember(ctx, updateInput)
}

func (r *MemberRepo) Checkout(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	user, err := r.GetMemberByUsername(ctx, username, constants.CART)
	if err != nil {
		return nil, 0, utils.LogError("err retrieving user", err)
	}

	if data.MemberTypes[user.Type] < len(movieIDs)+len(user.Checkedout) {
		return []string{"member limit exceeded"}, 0, nil
	}

	movies, err := r.movieRepo.GetMoviesByID(ctx, movieIDs, constants.NOT_CART)
	if err != nil {
		return nil, 0, utils.LogError("err retrieving movies", err)
	}

	return r.performCheckout(ctx, user, movies)
}

func (r *MemberRepo) Return(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	var messages []string
	var returned int

	movies, err := r.movieRepo.GetMoviesByID(ctx, movieIDs, constants.NOT_CART)
	if err != nil {
		return nil, 0, utils.LogError("err fetching movies for return", err)
	}

	name, err := attributevalue.Marshal(username)
	if err != nil {
		return nil, 0, utils.LogError("err marshalling username", err)
	}

	for _, movie := range movies {
		updateInput, err := r.getReturnInput(movie, name)
		if err != nil {
			messages = append(messages, utils.LogError(fmt.Sprintf("preparing return for %s", movie.Title), err).Error())
			continue
		}

		ok, err := r.updateMember(ctx, updateInput)
		if err != nil || !ok {
			messages = append(messages, utils.LogError(fmt.Sprintf("returning %s", movie.Title), err).Error())
			continue
		}

		ok, err = r.movieRepo.Return(ctx, movie)
		if err != nil || !ok {
			messages = append(messages, utils.LogError(fmt.Sprintf("updating inventory for %s", movie.Title), err).Error())
			continue
		}
		returned++
	}
	return messages, returned, nil
}

func (r *MemberRepo) SetMemberAPIChoice(ctx context.Context, username, apiChoice string) error {
	if apiChoice != constants.REST_API && apiChoice != constants.GRAPHQL_API {
		return utils.LogError(fmt.Sprintf("%s is not valid api selection. Choices: \"REST\" and \"GraphQL\"", apiChoice),
			fmt.Errorf("unexpected API Choice: %s", apiChoice))
	}
	name, err := attributevalue.Marshal(username)
	if err != nil {
		return utils.LogError(fmt.Sprintf("failed to marshal username %s: %v", username, err), err)
	}
	updateExpr := "SET api_choice = :api_choice"
	expressionAttrs := map[string]types.AttributeValue{
		":api_choice": &types.AttributeValueMemberS{Value: apiChoice},
	}
	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       map[string]types.AttributeValue{constants.USERNAME: name},
		ExpressionAttributeValues: expressionAttrs,
		UpdateExpression:          &updateExpr,
		ReturnValues:              types.ReturnValueUpdatedNew,
	}
	_, err = r.client.UpdateItem(ctx, updateInput)
	if err != nil {
		return utils.LogError(fmt.Sprintf("failed to update member %s api choice: %v", username, err), err)
	}
	return nil
}

func (r *MemberRepo) performCheckout(ctx context.Context, user data.Member, movies []data.Movie) ([]string, int, error) {
	var rented int
	var messages []string

	for _, movie := range movies {
		if movie.Inventory < 0 {
			messages = append(messages, fmt.Sprintf("%s is out of stock", movie.Title))
			continue
		}
		if utils.Contains(user.Checkedout, movie.ID) {
			messages = append(messages, fmt.Sprintf("%s is already checked out", movie.Title))
			continue
		}
		if !utils.Contains(user.Cart, movie.ID) {
			messages = append(messages, fmt.Sprintf("%s is not in cart", movie.Title))
			continue
		}
		if err := r.checkoutMovie(ctx, user, movie); err != nil {
			messages = append(messages, err.Error())
			continue
		}
		rented++
	}
	return messages, rented, nil
}

func (r *MemberRepo) checkoutMovie(ctx context.Context, user data.Member, movie data.Movie) error {
	ok, err := r.movieRepo.Rent(ctx, movie)
	if err != nil || !ok {
		return utils.LogError(fmt.Sprintf("renting %s", movie.Title), err)
	}

	ok, err = r.ModifyCart(ctx, user.Username, movie.ID, constants.DELETE, constants.CHECKOUT)
	if err != nil || !ok {
		r.movieRepo.Return(ctx, movie)
		return utils.LogError(fmt.Sprintf("removing %s from cart", movie.Title), err)
	}
	return nil
}

func (r *MemberRepo) updateMember(ctx context.Context, input *dynamodb.UpdateItemInput) (bool, error) {
	_, err := r.client.UpdateItem(ctx, input)
	if err != nil {
		return false, utils.LogError("updating member: ", err)
	}
	return true, nil
}

func (r *MemberRepo) getReturnInput(movie data.Movie, name types.AttributeValue) (*dynamodb.UpdateItemInput, error) {
	expr := "DELETE checked_out :checked_out ADD rented :rented"
	attrs := map[string]types.AttributeValue{
		":checked_out": &types.AttributeValueMemberSS{Value: []string{movie.ID}},
		":rented":      &types.AttributeValueMemberSS{Value: []string{movie.ID}},
	}
	return &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       map[string]types.AttributeValue{constants.USERNAME: name},
		ExpressionAttributeValues: attrs,
		UpdateExpression:          &expr,
		ReturnValues:              types.ReturnValueUpdatedNew,
	}, nil
}

func buildCartUpdateExpr(movieID, updateKey string, checkingOut bool) (string, map[string]types.AttributeValue) {
	expr := fmt.Sprintf("%s cart :cart", updateKey)
	attrs := map[string]types.AttributeValue{
		":cart": &types.AttributeValueMemberSS{Value: []string{movieID}},
	}
	if checkingOut {
		expr += " ADD checked_out :checked_out"
		attrs[":checked_out"] = &types.AttributeValueMemberSS{Value: []string{movieID}}
	}
	return expr, attrs
}

func (r *MemberRepo) GetIniitialVotingSlate(ctx context.Context) ([]string, error) {
	if r.centroids.Size() == 0 {
		return nil, utils.LogError("centroid cache failed to initialize; cannot support rec engine", nil)
	}

	usedIDs, attempts, errs := make(map[string]bool), 0, make([]error, 0)
	i := 0
	for i < constants.MAX_MOVIE_SUGGESTIONS && attempts < constants.MAX_MOVIE_SUGGESTIONS*3 {
		attempts++
		mid, err := r.centroidsToMovies.GetRandomMovieFromCentroid(r.randGen.Intn(r.centroids.Size()))
		if err != nil {
			errs = append(errs, err)
			utils.LogError("failed to attain random movie from centroid", err)
			continue
		}
		if usedIDs[mid] {
			continue
		}
		usedIDs[mid] = true
		i++
	}
	return utils.GetSliceFromMapKeys(usedIDs), errors.Join(errs...)
}

func (r *MemberRepo) IterateRecommendationVoting(
	ctx context.Context,
	currentMood data.MovieMetrics,
	iteration, numPrevSelected int,
	movieIDs []string) (data.MovieMetrics, []string, error) {

	updatedMood, err := r.UpdateMood(ctx, currentMood, numPrevSelected, movieIDs)
	if err != nil {
		return data.MovieMetrics{}, nil, utils.LogError("updating mood", err)
	}
	newCentroids, err := r.centroids.GetKNearestCentroidsFromMood(updatedMood, constants.MAX_CENTROIDS_COUNT-iteration)
	if err != nil {
		return data.MovieMetrics{}, nil, utils.LogError("getting new centroids", err)
	}

	movieRecs, originalMovies := make(map[string]bool), make(map[string]bool)
	var errs []error
	if len(newCentroids) > 0 {
		for _, mid := range movieIDs {
			originalMovies[mid] = true
		}
		for i := 0; i < constants.MAX_MOVIE_SUGGESTIONS-(iteration*2); i++ {
			centroid := newCentroids[rand.Intn(len(newCentroids))]
			movieID, attempts := "", 0
			for attempts < 10 {
				attempts++
				movieID, err = r.centroidsToMovies.GetRandomMovieFromCentroid(centroid)
				if err != nil {
					errs = append(errs, utils.LogError("error getting random movie from centroid", err))
					continue
				}
				if !originalMovies[movieID] && !movieRecs[movieID] {
					break // found a unique movie
				}
			}
			if movieID == "" {
				errs = append(errs, utils.LogError(fmt.Sprintf("failed to find random movie for centroid %v", centroid), nil))
				continue
			}
			movieRecs[movieID] = true
		}
	}
	return updatedMood, utils.GetSliceFromMapKeys(movieRecs), errors.Join(errs...)
}

func (r *MemberRepo) UpdateMood(ctx context.Context, currentMood data.MovieMetrics, numPrevSelected int, movieIDs []string) (data.MovieMetrics, error) {
	accMood, updateCount, errs := utils.AccumulateMovieMetricsWithWeight(data.MovieMetrics{}, currentMood, numPrevSelected), 0, []error{}
	for _, mid := range movieIDs {
		metrics, err := r.movieRepo.GetMovieMetrics(ctx, mid)
		if err != nil {
			errs = append(errs, utils.LogError(fmt.Sprintf("failed to retrieve movie metrics for %s", mid), err))
			continue
		}
		updateCount++
		accMood = utils.AccumulateMovieMetricsWithWeight(accMood, metrics, 1)
	}
	return utils.AverageMetrics(accMood, numPrevSelected+updateCount), errors.Join(errs...)
}

func (r *MemberRepo) GetVotingFinalPicks(ctx context.Context, mood data.MovieMetrics) ([]string, error) {
	centroidIDs, err := r.centroids.GetKNearestCentroidsFromMood(mood, constants.NUMBER_FINAL_PICKS)
	if err != nil {
		return nil, utils.LogError("failed to get centroid neighbors", err)
	}

	suggestions := make([]string, 0, len(centroidIDs))
	for _, id_ := range centroidIDs {
		mid, err := r.getNearestNeighborInCentroid(ctx, id_, mood)
		if err != nil {
			utils.LogError(fmt.Sprintf("failed to come up with neareset movie in centroid %v", id_), nil)
		}
		suggestions = append(suggestions, mid)
	}
	return suggestions, nil
}

func (r *MemberRepo) getNearestNeighborInCentroid(ctx context.Context, centroidID int, mood data.MovieMetrics) (string, error) {
	movieMetrics, err := r.getMovieMetricsForCentroid(ctx, centroidID)
	if err != nil {
		return "", err
	}

	minDistance := math.MaxFloat64
	nearestNeighbor := ""
	for mid, mets := range movieMetrics {
		d := utils.MetricDistance(mood, mets)
		if d < minDistance {
			nearestNeighbor = mid
			minDistance = d
		}
	}
	return nearestNeighbor, nil
}

func (r *MemberRepo) getMovieMetricsForCentroid(ctx context.Context, centroidID int) (map[string]data.MovieMetrics, error) {
	movies, err := r.centroidsToMovies.GetMovieIDsByCentroid(centroidID)
	if err != nil || len(movies) == 0 {
		return nil, utils.LogError(fmt.Sprintf("no movies found for centroid id %v", centroidID), nil)
	}

	mets := make(map[string]data.MovieMetrics)
	for _, mid := range movies {
		// check cache
		if metrics, ok := r.movieMetricCache[mid]; ok {
			mets[mid] = metrics
			continue
		}

		// populate cache
		metrics, err := r.movieRepo.GetMovieMetrics(ctx, mid)
		if err != nil {
			utils.LogError(fmt.Sprintf("failed to get movie id for movie id %s", mid), nil)
			continue
		}
		r.movieMetricCache[mid] = metrics
		mets[mid] = metrics
	}
	return mets, nil
}
