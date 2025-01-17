namespace DotnetApi.Queries
{
    public class PropertyQuery : ListQuery
    {
        public string LocationId { get; set; }

        public string DistrictId { get; set; }

        public string PowerLineId { get; set; }

        public string OwnerId { get; set; }
    }
}
