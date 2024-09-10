package db

import (
	// "blockbuster/api/endpoints"
	// "blockbuster/api/endpoints"

	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const membersTableName = "BluckBoster_members"

type MemberRepo struct {
	svc       dynamodb.DynamoDB
	tableName string
}

func NewMembersRepo() MemberRepo {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return MemberRepo{
		svc:       *dynamodb.New(sess),
		tableName: membersTableName,
	}
}

func (r MemberRepo) GetMemberByUsername(username string) (bool, Member, error) {
	// @TODO: never actually return error (it's always nil)
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			USERNAME: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(username),
					},
				},
			},
		},
	}

	result, err := r.svc.Query(queryInput)
	if err != nil {
		log.Fatalf("Got error calling svc.Query %s\n", err)
	}
	if len(result.Items) == 0 {
		log.Printf("Could not find member with username: %s\n", username)
		return false, Member{}, nil
	}
	member := Member{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &member)
	if err != nil {
		log.Fatalf("Failed to unmarshall data %s\n", err)
	}
	return true, member, nil
}

func (r MemberRepo) AddToCart(username, movieID string) (
	bool, *dynamodb.UpdateItemOutput, error,
) {
	ids, err := r.GetCartIDs(username)
	if err != nil {
		log.Printf("Failed checking to see if %s is in %s cart", movieID, username)
		return false, nil, err
	}
	for _, id := range ids {
		if id == movieID {
			return false, nil, nil
		}
	}

	var c []*dynamodb.AttributeValue
	c = append(c, &dynamodb.AttributeValue{S: aws.String(movieID)})
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			USERNAME: {
				S: aws.String(username),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":cart": {
				L: c,
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET cart = list_append(if_not_exists(cart, :empty_list), :cart)"),
	}
	response, err := r.svc.UpdateItem(updateInput)
	if err != nil {
		log.Printf("Failed to add movie %s to %s cart\n %s\n", movieID, username, err)
		return false, response, err
	}
	return true, response, nil
}

func (r MemberRepo) GetCartIDs(username string) ([]string, error) {
	input := &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			USERNAME: {
				S: aws.String(username),
			},
		},
		ProjectionExpression: aws.String("cart"),
	}

	response, err := r.svc.GetItem(input)
	if err != nil {
		log.Printf("Err fetching movies from cloud: %s\n", err)
		return nil, err
	}

	cart := Cart{}
	err = dynamodbattribute.UnmarshalMap(response.Item, &cart)
	if err != nil {
		log.Printf("Err unmarshalling movies from: %s\n", err)
		return nil, err
	}
	return cart.Cart, err
}
