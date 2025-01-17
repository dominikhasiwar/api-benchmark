package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DbContext struct {
	client    *dynamodb.Client
	tableName string
}

// var (
// 	instance *DbContext
// 	once     sync.Once
// )

func InitDbContext(ctx context.Context, region string, tableName string, svcEndpoint string) (*DbContext, error) {
	cfg, err := config.LoadDefaultConfig(ctx, func(options *config.LoadOptions) error {
		options.Region = region
		options.RetryMaxAttempts = 5
		return nil
	})
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(options *dynamodb.Options) {
		if svcEndpoint != "" {
			options.BaseEndpoint = aws.String(svcEndpoint)
			options.Credentials = credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")
		}
	})

	return &DbContext{
		client:    client,
		tableName: tableName,
	}, nil
}

func (db *DbContext) GetList(ctx context.Context, filter expression.ConditionBuilder, maxResults int32, lastEvaluatedKey string) ([]map[string]types.AttributeValue, string, error) {
	var items []map[string]types.AttributeValue
	var lastEvaluatedKeyMap map[string]types.AttributeValue

	scanLimit := aws.Int32(500)

	if lastEvaluatedKey != "" {
		lastEvaluatedKeyMap = map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: lastEvaluatedKey},
		}
	}

	input := &dynamodb.ScanInput{
		TableName:         aws.String(db.tableName),
		ExclusiveStartKey: lastEvaluatedKeyMap,
		Limit:             scanLimit,
	}

	if filter.IsSet() {
		condition, err := expression.NewBuilder().WithFilter(filter).Build()
		if err != nil {
			return nil, "", err
		}

		input.ExpressionAttributeNames = condition.Names()
		input.ExpressionAttributeValues = condition.Values()
		input.FilterExpression = condition.Filter()
	}

	for {
		input.ExclusiveStartKey = lastEvaluatedKeyMap

		output, err := db.client.Scan(ctx, input)
		if err != nil {
			return nil, "", err
		}

		items = append(items, output.Items...)

		if output.LastEvaluatedKey == nil {
			lastEvaluatedKey = ""
			break
		}

		lastEvaluatedKeyMap = output.LastEvaluatedKey

		if maxResults > 0 && len(items) >= int(maxResults) {
			lastEvaluatedKey = lastEvaluatedKeyMap["Id"].(*types.AttributeValueMemberS).Value
			break
		}
	}

	return items, lastEvaluatedKey, nil
}

func (db *DbContext) GetSingle(ctx context.Context, filter expression.ConditionBuilder) (map[string]types.AttributeValue, error) {
	entities, _, err := db.GetList(ctx, filter, 1, "")
	if err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return nil, nil
	}

	return entities[0], nil
}

func (db *DbContext) Count(ctx context.Context, filter expression.ConditionBuilder) (int64, error) {
	var lastEvaluatedKeyMap map[string]types.AttributeValue
	var count int64

	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return 0, err
	}

	for {
		input := &dynamodb.ScanInput{
			TableName:                 aws.String(db.tableName),
			ExpressionAttributeNames:  condition.Names(),
			ExpressionAttributeValues: condition.Values(),
			FilterExpression:          condition.Filter(),
			Select:                    types.SelectCount,
			ExclusiveStartKey:         lastEvaluatedKeyMap,
		}

		output, err := db.client.Scan(context.Background(), input)
		if err != nil {
			return 0, err
		}

		lastEvaluatedKeyMap = output.LastEvaluatedKey
		count += int64(output.Count)

		if lastEvaluatedKeyMap == nil {
			break
		}
	}

	return count, nil
}

func (db *DbContext) Save(ctx context.Context, entity interface{}) (*dynamodb.PutItemOutput, error) {
	entityParsed, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.tableName),
		Item:      entityParsed,
	}

	return db.client.PutItem(ctx, input)
}

func (db *DbContext) SaveBatch(ctx context.Context, entities []interface{}) error {
	writeRequests := make([]types.WriteRequest, len(entities))
	for i, entity := range entities {
		item, err := attributevalue.MarshalMap(entity)
		if err != nil {
			return err
		}

		writeRequests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		}
	}

	batchSize := 25
	start := 0
	end := start + batchSize

	for start < len(writeRequests) {
		if end > len(writeRequests) {
			end = len(writeRequests)
		}

		// Prepare batch input
		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				db.tableName: writeRequests[start:end],
			},
		}

		_, err := db.client.BatchWriteItem(ctx, input)
		if err != nil {
			return err
		}

		start = end
		end += batchSize

	}

	return nil
}

func (db *DbContext) Delete(ctx context.Context, condition map[string]interface{}) (*dynamodb.DeleteItemOutput, error) {
	conditionParsed, err := attributevalue.MarshalMap(condition)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(db.tableName),
		Key:       conditionParsed,
	}

	return db.client.DeleteItem(ctx, input)
}
