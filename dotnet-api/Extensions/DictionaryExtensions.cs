using System.Reflection;
using Amazon.DynamoDBv2.Model;

namespace DotnetApi.Extensions
{
  public static class DictionaryExtensions
  {
    public static T ConvertToModel<T>(this Dictionary<string, AttributeValue> item)
    {
      var model = Activator.CreateInstance<T>();
      var properties = typeof(T).GetProperties(BindingFlags.Public | BindingFlags.Instance);

      foreach (var property in properties)
      {
        if (item.TryGetValue(property.Name, out var value))
        {
          if (property.PropertyType == typeof(DateTime) && DateTime.TryParse(value.S, out var dateTimeValue))
          {
            property.SetValue(model, dateTimeValue);
          }
          else
          {
            property.SetValue(model, value.S);
          }
        }
      }

      return model;
    }

    public static Dictionary<string, AttributeValue> ModelToDictionary<T>(this T model)
    {
      var dictionary = new Dictionary<string, AttributeValue>();
      var properties = typeof(T).GetProperties(BindingFlags.Public | BindingFlags.Instance);

      foreach (var property in properties)
      {
        var value = property.GetValue(model)?.ToString();
        if (value != null)
        {
          dictionary[property.Name] = new AttributeValue { S = value };
        }
      }

      return dictionary;
    }
  }
}
