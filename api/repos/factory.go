package repos

import (
	utils "blockbuster/api/utils"
)

// GetMovieRepoWithDynamo returns a new MovieRepo using the default DynamoDB client.
func NewMovieRepoWithDynamo() *MovieRepo {
	client := utils.GetDynamoClient()
	return NewMovieRepo(client)
}

// GetMemberRepoWithDynamo returns a new MemberRepo with its own MovieRepo dependency.
func NewMemberRepoWithDynamo() *MemberRepo {
	client := utils.GetDynamoClient()
	movieRepo := NewMovieRepo(client)
	return NewMembersRepo(client, movieRepo)
}
