using DotnetApi.Models;
using DotnetApi.Queries;
using DotnetApi.Repositories;
using Microsoft.AspNetCore.Mvc;

namespace DotnetApi.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class UserController : ControllerBase
    {
        private readonly IUserRepository _userRepository;

        public UserController(IUserRepository userRepository)
        {
            _userRepository = userRepository;
        }

        [HttpGet]
        public async Task<ListResponseModel<UserModel>> GetUsers(ListQuery query = null)
        {
            return await _userRepository.GetUsers(query);
        }

        [HttpGet("{id}")]
        public async Task<UserModel> GetUser(string id)
        {
            return await _userRepository.GetUser(id);
        }

        [HttpPost]
        public async Task<UserModel> CreateUser(SaveUserModel model)
        {
            return await _userRepository.CreateUser(model);
        }

        [HttpPut("{id}")]
        public async Task<UserModel> UpdateUser(string id, SaveUserModel model)
        {
            return await _userRepository.UpdateUser(id, model);
        }

        [HttpDelete("{id}")]
        public async Task DeleteUser(string id)
        {
            await _userRepository.DeleteUser(id);
        }

        [HttpPost("import")]
        public async Task<ImportResultModel> ImportUsers(IFormFile file, string password = null)
        {
            return await _userRepository.ImportUsers(file.OpenReadStream(), password);
        }

        [HttpDelete]
        public async Task ClearUsers()
        {
            await _userRepository.ClearUsers();
        }
    }
}
