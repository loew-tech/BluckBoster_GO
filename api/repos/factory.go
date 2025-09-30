package repos

import (
	"sync"

	"blockbuster/api/api_cache"
	"blockbuster/api/utils"
)

var (
	movieRepoOnce  sync.Once
	memberRepoOnce sync.Once

	movieRepoInstance  ReadWriteMovieRepo
	memberRepoInstance MemberRepoInterface
)

// NewMovieRepoWithDynamo returns a singleton MovieRepo using a shared DynamoDB client.
func NewMovieRepoWithDynamo() ReadWriteMovieRepo {
	movieRepoOnce.Do(func() {
		client := utils.GetDynamoClient()
		movieRepoInstance = NewDynamoMovieRepo(client)
	})
	return movieRepoInstance
}

// NewMemberRepoWithDynamo returns a singleton MemberRepo using a shared MovieRepo.
func NewMemberRepoWithDynamo() MemberRepoInterface {
	memberRepoOnce.Do(func() {
		client := utils.GetDynamoClient()
		movieRepo := NewMovieRepoWithDynamo()
		memberRepoInstance = NewMembersRepo(client, movieRepo, api_cache.GetDynamoClientCentroidCache(), api_cache.InitCentroidsToMoviesCache(movieRepo.GetMoviesByPage))
	})
	return memberRepoInstance
}
