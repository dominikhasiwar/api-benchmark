namespace DotnetApi.Queries
{
    public class DispatchQuery : ListQuery
    {
        public string LocationId { get; set; }

        public string DistrictId { get; set; }

        public string PowerLineId { get; set; }

        public string ContactPersonId { get; set; }

        public string OwnerId { get; set; }
    }
}
