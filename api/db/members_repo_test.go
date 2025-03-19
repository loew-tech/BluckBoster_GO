package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestRepo() MemberRepo {
	mockMember, mockMovie := MockDynamoMemberClient{}, MockMovieDynamoClient{}
	return NewMembersRepo(mockMember, NewMovieRepo(mockMovie))
}

func TestGetMemberByUsername(t *testing.T) {

	testMemberRepo := getTestRepo()
	found, user, err := testMemberRepo.GetMemberByUsername(TestMember.Username, false)

	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, TestMember, user)
}

func TestGetMemberByUsernameReturnNil(t *testing.T) {

	testMemberRepo := getTestRepo()
	found, user, err := testMemberRepo.GetMemberByUsername(TestMember.Username, true)

	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, TestMember.Username, user.Username)
	assert.Equal(t, TestMember.Cart, user.Cart)
	assert.Equal(t, TestMember.Checkedout, user.Checkedout)
	assert.Equal(t, TestMember.Type, user.Type)
	assert.NotEqual(t, TestMember, user)
}
