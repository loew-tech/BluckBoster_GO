package repos

import (
	utils "blockbuster/api/utils"
)

// GetMovieRepoWithDynamo returns a new MovieRepo using the default DynamoDB client.
func NewMovieRepoWithDynamo() ReadWriteMovieRepo {
	client := utils.GetDynamoClient()
	return NewDynamoMovieRepo(client)
}

// GetMemberRepoWithDynamo returns a new MemberRepo with its own MovieRepo dependency.
func NewMemberRepoWithDynamo() *MemberRepo {
	client := utils.GetDynamoClient()
	movieRepo := NewDynamoMovieRepo(client)
	return NewMembersRepo(client, movieRepo)
}
