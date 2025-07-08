package repos

import (
	"blockbuster/api/utils"
)

func NewMovieRepoWithDynamo() ReadWriteMovieRepo {
	client := utils.GetDynamoClient()
	return newDynamoMovieRepo(client)
}

func NewMemberRepoWithDynamo() MemberRepoInterface {
	client := utils.GetDynamoClient()
	movieRepo := newDynamoMovieRepo(client)
	return newMembersRepo(client, movieRepo)
}
