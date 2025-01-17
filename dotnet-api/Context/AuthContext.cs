using System.Security.Claims;
using DotnetApi.Models;

namespace DotnetApi.Context
{
    public interface IAuthContext
    {
        UserInfoModel GetCurrentUser();
    }

    public class AuthContext : IAuthContext
    {
        private readonly IHttpContextAccessor _contextAccessor;

        public AuthContext(IHttpContextAccessor contextAccessor)
        {
            _contextAccessor = contextAccessor;
        }

        public UserInfoModel GetCurrentUser()
        {
            var principal = _contextAccessor.HttpContext.User;

            var user = new UserInfoModel();

            if (principal.Identity.IsAuthenticated)
            {
                user.UserName = principal.FindFirstValue("preferred_username");
            }

            return user;
        }
    }
}