using AutoMapper;
using DotnetApi.Entities;
using DotnetApi.Models;

namespace DotnetApi.Profiles
{
    public class UserProfile : Profile
    {
        public UserProfile()
        {
            CreateMap<User, UserModel>();

            CreateMap<SaveUserModel, User>();
        }
    }
}
