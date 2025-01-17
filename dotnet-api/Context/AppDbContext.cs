using Amazon.DynamoDBv2;
using Amazon.DynamoDBv2.DataModel;
using Amazon.DynamoDBv2.Model;
using DotnetApi.Entities;
using DotnetApi.Extensions;
using DotnetApi.Helpers;

namespace DotnetApi
{
    public interface IAppDbContext
    {
        Task<(string LastEvaluatedKey, List<T> List)> GetList<T>(string lastEvaluatedKey = null, IFilterCondition filterCondition = null, int maxItemsCount = 0);

        Task<T> GetSingle<T>(string id);

        Task<T> GetSingle<T>(IFilterCondition filterCondition);

        Task Save<T>(T model) where T : EntityBase;

        Task SaveBatch<T>(List<T> list);

        Task Delete<T>(string id);

        Task Delete<T>(IFilterCondition filterCondition);
    }

    public class AppDbContext : IAppDbContext
    {
        private const int SCAN_LIMIT = 200;

        private readonly string _tableName;
        private readonly IDynamoDBContext _dynamoDBContext;
        private readonly IAmazonDynamoDB _dbClient;
        private readonly DynamoDBOperationConfig _operationConfig;

        public AppDbContext(string tableName, IDynamoDBContext dynamoDBContext, IAmazonDynamoDB dbClient)
        {
            _tableName = tableName;
            _dynamoDBContext = dynamoDBContext;
            _dbClient = dbClient;
            _operationConfig = new DynamoDBOperationConfig
            {
                OverrideTableName = _tableName
            };
        }

        public async Task<(string LastEvaluatedKey, List<T> List)> GetList<T>(string lastEvaluatedKey = null, IFilterCondition filterCondition = null, int maxItemsCount = 0)
        {
            var items = new List<T>();
            var lastEvaluatedKeyMap = new Dictionary<string, AttributeValue>();

            if (!string.IsNullOrEmpty(lastEvaluatedKey))
            {
                lastEvaluatedKeyMap.Add(nameof(EntityBase.Id), new AttributeValue { S = lastEvaluatedKey });
            }

            do
            {
                var request = new ScanRequest
                {
                    TableName = _tableName,
                    FilterExpression = filterCondition?.GetFilterExpression(),
                    ExpressionAttributeValues = filterCondition?.GetExpressionValues(),
                    ExclusiveStartKey = lastEvaluatedKeyMap,
                    Limit = SCAN_LIMIT
                };

                var response = await _dbClient.ScanAsync(request);
                lastEvaluatedKeyMap = response.LastEvaluatedKey;

                foreach (var item in response.Items)
                {
                    var model = item.ConvertToModel<T>();

                    items.Add(model);
                }

                if (lastEvaluatedKeyMap.Count == 0)
                {
                    break;
                }

                if (maxItemsCount > 0 && items.Count >= maxItemsCount)
                {
                    break;
                }
            }
            while (lastEvaluatedKeyMap.Any());

            if (lastEvaluatedKeyMap.ContainsKey(nameof(EntityBase.Id)))
            {
                lastEvaluatedKey = lastEvaluatedKeyMap[nameof(EntityBase.Id)].S;
            }

            return (lastEvaluatedKey, items);
        }

        public async Task<T> GetSingle<T>(string id)
        {
            var item = await _dynamoDBContext.LoadAsync<T>(id, _operationConfig);

            return item;
        }

        public async Task<T> GetSingle<T>(IFilterCondition filterCondition)
        {
            var response = await GetList<T>(filterCondition: filterCondition, maxItemsCount: 1);

            return response.List.FirstOrDefault();
        }

        public Task Save<T>(T item) where T : EntityBase
        {
            item.Id ??= Guid.NewGuid().ToString();

            return _dynamoDBContext.SaveAsync(item, _operationConfig);
        }

        public Task SaveBatch<T>(List<T> list)
        {
            var batchWrite = _dynamoDBContext.CreateBatchWrite<T>(_operationConfig);

            foreach (var item in list)
            {
                batchWrite.AddPutItem(item);
            }

            return batchWrite.ExecuteAsync();
        }

        public async Task Delete<T>(string id)
        {
            var item = await _dynamoDBContext.LoadAsync<T>(id, _operationConfig);

            if (item != null)
            {
                await _dynamoDBContext.DeleteAsync(item, _operationConfig);
            }
        }

        public async Task Delete<T>(IFilterCondition filterCondition)
        {
            var response = await GetList<T>(filterCondition: filterCondition);

            foreach (var item in response.List)
            {
                await _dynamoDBContext.DeleteAsync(item, _operationConfig);
            }
        }
    }
}