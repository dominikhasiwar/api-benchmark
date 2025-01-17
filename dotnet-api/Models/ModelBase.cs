namespace DotnetApi.Models
{
    public class ModelBase
    {
        public string Id { get; set; }

        public string Creator { get; set; }

        public DateTime Created { get; set; }

        public string Modifier { get; set; }

        public DateTime Modified { get; set; }
    }
}
