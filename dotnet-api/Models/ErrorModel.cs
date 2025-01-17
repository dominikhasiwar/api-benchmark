using System.Net;
using System.Text.Json;
using System.Xml;

namespace DotnetApi.Models
{
    public class ErrorModel
    {
        internal HttpStatusCode StatusCode { get; set; }

        public string ErrorMessage { get; set; }

        public override string ToString()
        {
            return JsonSerializer.Serialize(this);
        }
    }
}
