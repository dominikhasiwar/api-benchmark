namespace DotnetApi.Models
{
    public class UserModel : ModelBase
    {
        public string UserName { get; set; }

        public string FirstName { get; set; }

        public string LastName { get; set; }

        public string Fullname
        {
            get { return $"{FirstName} {LastName}".Trim(); }
        }

        public string Street { get; set; }

        public string City { get; set; }

        public string Zip { get; set; }

        public string Country { get; set; }

        public string Email { get; set; }

        public string PhoneNumber { get; set; }
    }
}
