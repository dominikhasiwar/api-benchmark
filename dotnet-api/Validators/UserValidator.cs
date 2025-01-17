using DotnetApi.Entities;
using DotnetApi.Helpers;
using DotnetApi.Models;
using FluentValidation;

namespace DotnetApi.Validators
{
    public class UserValidator : AbstractValidator<SaveUserModel>
    {
        private readonly IAppDbContext _dbContext;

        public UserValidator(IAppDbContext dbContext)
        {
            _dbContext = dbContext;

            RuleFor(x => x.UserName).NotEmpty().Must(UniqueUserName).WithMessage("User name already exists");

            RuleFor(x => x.FirstName).NotEmpty();

            RuleFor(x => x.LastName).NotEmpty();

            RuleFor(x => x.Email).NotEmpty().EmailAddress();
        }

        private bool UniqueUserName(string value)
        {
            var filterCondition = FilterCondition.For(nameof(User.UserName)).Equals(value);

            var user = _dbContext.GetSingle<User>(filterCondition).Result;

            return user == null;
        }
    }
}
