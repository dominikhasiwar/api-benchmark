using Microsoft.AspNetCore.Mvc;

namespace DotnetApi.Queries
{
    public class ListQuery
    {
        [FromQuery]
        public string LastEvaluatedKey { get; set; }

        [FromQuery]
        public string TextQuery { get; set; }
    }

}