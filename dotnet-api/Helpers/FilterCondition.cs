using System.Text;
using Amazon.DynamoDBv2.Model;
using Microsoft.AspNetCore.Mvc.ModelBinding;

namespace DotnetApi.Helpers
{
    public enum FilterOperator
    {
        Equals,
        NotEquals,
        Contains,
    }

    public interface IFilterCondition
    {
        string GetFilterExpression();

        Dictionary<string, AttributeValue> GetExpressionValues();
    }

    public class FilterCondition : IFilterCondition
    {
        private static readonly Random _random = new Random();
        private const string chars = "abcdefghijklmnopqrstuvwxyz0123456789";

        private readonly string _expression;

        private readonly Dictionary<string, AttributeValue> _values = new Dictionary<string, AttributeValue>();

        public FilterCondition(string propertyName, string value, FilterOperator filterOperator)
        {
            var valueId = $":{GenerateRandomString()}";

            switch (filterOperator)
            {
                case FilterOperator.Equals:
                    _expression = $"{propertyName} = {valueId}";
                    break;
                case FilterOperator.NotEquals:
                    _expression = $"{propertyName} <> {valueId}";
                    break;
                case FilterOperator.Contains:
                    _expression = $"contains({propertyName}, {valueId})";
                    break;
                default:
                    throw new ArgumentOutOfRangeException(nameof(filterOperator), filterOperator, null);
            }

            _values.Add(valueId, new AttributeValue { S = value });
        }

        public FilterCondition(string expression, Dictionary<string, AttributeValue> values)
        {
            _expression = expression;
            _values = values;
        }

        public static FilterCondition Empty()
        {
            return new FilterCondition(null, new Dictionary<string, AttributeValue>());
        }

        public static FilterExpressionName For(string propertyName)
        {
            return new FilterExpressionName(propertyName);
        }

        public FilterCondition And(FilterCondition right)
        {
            var values = _values.Concat(right.GetExpressionValues()).ToDictionary(kvp => kvp.Key, kvp => kvp.Value);

            if (string.IsNullOrEmpty(_expression))
            {
                return right;
            }
            else
            {
                return new FilterCondition($"({_expression} AND {right.GetFilterExpression()})", values);
            }
        }

        public FilterCondition Or(FilterCondition right)
        {
            var values = _values.Concat(right.GetExpressionValues()).ToDictionary(kvp => kvp.Key, kvp => kvp.Value);

            if (string.IsNullOrEmpty(_expression))
            {
                return right;
            }
            else
            {
                return new FilterCondition($"({_expression} OR {right.GetFilterExpression()})", values);
            }
        }

        private static string GenerateRandomString()
        {
            var stringBuilder = new StringBuilder(5);
            for (int i = 0; i < 5; i++)
            {
                stringBuilder.Append(chars[_random.Next(chars.Length)]);
            }
            return stringBuilder.ToString();
        }

        public string GetFilterExpression()
        {
            return _expression;
        }

        public Dictionary<string, AttributeValue> GetExpressionValues()
        {
            return _values;
        }
    }

    public class FilterExpressionName
    {
        public string PropertyName { get; }

        public FilterExpressionName(string propertyName)
        {
            PropertyName = propertyName;
        }

        public FilterCondition Equals(string value)
        {
            return new FilterCondition(PropertyName, value, FilterOperator.Equals);
        }

        public FilterCondition Contains(string value)
        {
            return new FilterCondition(PropertyName, value, FilterOperator.Contains);
        }
    }
}