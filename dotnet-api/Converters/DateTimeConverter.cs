using Amazon.DynamoDBv2.DataModel;
using Amazon.DynamoDBv2.DocumentModel;

namespace DotnetApi.Converters
{
    public class DateTimeConverter : IPropertyConverter
    {
        public object FromEntry(DynamoDBEntry entry)
        {
            if (entry == null || string.IsNullOrEmpty(entry.AsString()))
            {
                return DateTime.MinValue;
            }
            return DateTime.Parse(entry.AsString());
        }

        public DynamoDBEntry ToEntry(object value)
        {
            if (value == null)
            {
                return null;
            }
            return ((DateTime)value).ToString("o"); // ISO 8601 format
        }
    }

}