using Amazon.DynamoDBv2.DataModel;
using DotnetApi.Converters;

namespace DotnetApi.Entities
{
    public abstract class EntityBase
    {
        [DynamoDBHashKey]
        public string Id { get; set; }

        [DynamoDBProperty]
        public string Creator { get; set; }

        [DynamoDBProperty(Converter = typeof(DateTimeConverter))]
        public DateTime Created { get; set; }

        [DynamoDBProperty]
        public string Modifier { get; set; }

        [DynamoDBProperty(Converter = typeof(DateTimeConverter))]
        public DateTime Modified { get; set; }
    }
}
