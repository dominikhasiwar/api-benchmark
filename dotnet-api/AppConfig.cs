namespace DotnetApi
{
    public interface IAppConfig
    {
        DynamoDbConfig DynamoDb { get; }

        AzureAdConfig AzureAd { get; }
    }

    public class AppConfig : IAppConfig
    {
        public DynamoDbConfig DynamoDb { get; set; }

        public AzureAdConfig AzureAd { get; set; }
    }

    public class DynamoDbConfig
    {
        public string TableName { get; set; }

        public string ServiceUrl { get; set; }
    }

    public class AzureAdConfig
    {
        public string ClientId { get; set; }

        public string TenantId { get; set; }

        public string Scopes { get; set; }
    }
}
