using System.Text.Json.Serialization;
using Amazon.Lambda.APIGatewayEvents;
using DotnetApi.Models;

namespace DotnetApi
{
    [JsonSerializable(typeof(APIGatewayHttpApiV2ProxyRequest))]
    [JsonSerializable(typeof(APIGatewayHttpApiV2ProxyResponse))]
    [JsonSerializable(typeof(string))]
    [JsonSerializable(typeof(bool))]
    [JsonSerializable(typeof(int))]
    [JsonSerializable(typeof(List<string>))]
    [JsonSerializable(typeof(Dictionary<string, string>))]
    [JsonSerializable(typeof(ErrorModel))]
    [JsonSerializable(typeof(ModelBase))]
    [JsonSerializable(typeof(UserModel))]
    [JsonSerializable(typeof(UserModel[]))]
    [JsonSerializable(typeof(SaveUserModel))]
    [JsonSourceGenerationOptions(WriteIndented = true)]
    public partial class ApiSerializerContext : JsonSerializerContext
    {
    }
}
