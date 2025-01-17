namespace DotnetApi.Models
{
    public class ListResponseModel<T>
    {
        public string LastEvaluatedKey { get; set; }

        public T[] List { get; set; }
    }
}
