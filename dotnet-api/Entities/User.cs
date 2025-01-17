using Amazon.DynamoDBv2.DataModel;

namespace DotnetApi.Entities
{
    public class User : EntityBase
    {
        [DynamoDBProperty]
        public string UserName { get; set; }

        [DynamoDBProperty]
        public string FirstName { get; set; }

        [DynamoDBProperty]
        public string LastName { get; set; }

        [DynamoDBProperty]
        public string Street { get; set; }

        [DynamoDBProperty]
        public string City { get; set; }

        [DynamoDBProperty]
        public string Zip { get; set; }

        [DynamoDBProperty]
        public string Country { get; set; }

        [DynamoDBProperty]
        public string Email { get; set; }

        [DynamoDBProperty]
        public string PhoneNumber { get; set; }
    }
}
