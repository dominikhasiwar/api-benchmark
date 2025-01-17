using System.Data;
using System.Text;
using AutoMapper;
using DotnetApi.Context;
using DotnetApi.Entities;
using DotnetApi.Exceptions;
using DotnetApi.Models;
using DotnetApi.Queries;
using ExcelDataReader;

namespace DotnetApi.Repositories
{
    public interface IUserRepository
    {
        Task<ListResponseModel<UserModel>> GetUsers(ListQuery query = null);

        Task<UserModel> GetUser(string id);

        Task<UserModel> CreateUser(SaveUserModel model);

        Task<UserModel> UpdateUser(string id, SaveUserModel model);

        Task DeleteUser(string id);

        Task<ImportResultModel> ImportUsers(Stream stream, string password = null);

        Task ClearUsers();
    }

    public class UserRepository : IUserRepository
    {
        private readonly IAppDbContext _dbContext;
        private readonly IAuthContext _authContext;
        private readonly IMapper _mapper;

        public UserRepository(IAppDbContext dbContext, IAuthContext authContext, IMapper mapper)
        {
            Encoding.RegisterProvider(CodePagesEncodingProvider.Instance);

            _dbContext = dbContext;
            _authContext = authContext;
            _mapper = mapper;
        }

        public async Task<ListResponseModel<UserModel>> GetUsers(ListQuery query = null)
        {
            var result = await _dbContext.GetList<User>(query?.LastEvaluatedKey, null, 100);

            return new ListResponseModel<UserModel>
            {
                LastEvaluatedKey = result.LastEvaluatedKey,
                List = _mapper.Map<UserModel[]>(result.List)
            };
        }

        public async Task<UserModel> GetUser(string id)
        {
            var user = await _dbContext.GetSingle<User>(id) ?? throw new NotFoundException($"User with id {id} not found");

            return _mapper.Map<UserModel>(user);
        }

        public async Task<UserModel> CreateUser(SaveUserModel model)
        {
            var user = _mapper.Map<User>(model);

            user.Id = Guid.NewGuid().ToString();
            user.Creator = user.Modifier = _authContext.GetCurrentUser()?.UserName;
            user.Created = user.Modified = DateTime.UtcNow;

            await _dbContext.Save(user);

            return await GetUser(user.Id);
        }

        public async Task<UserModel> UpdateUser(string id, SaveUserModel model)
        {
            var user = await _dbContext.GetSingle<User>(id) ?? throw new NotFoundException($"User with id {id} not found");

            _mapper.Map(model, user);

            user.Modifier = _authContext.GetCurrentUser()?.UserName;
            user.Modified = DateTime.UtcNow;

            await _dbContext.Save(user);

            return await GetUser(user.Id);
        }

        public async Task DeleteUser(string id)
        {
            await _dbContext.Delete<User>(id);
        }

        public async Task<ImportResultModel> ImportUsers(Stream stream, string password = null)
        {
            var result = await _dbContext.GetList<User>();
            var existingPersons = result.List;
            var newUsers = new List<User>();

            using var reader = ExcelReaderFactory.CreateReader(stream, new ExcelReaderConfiguration
            {
                Password = password
            });

            var dataSet = reader.AsDataSet();

            foreach (var row in dataSet.Tables[0].Rows.OfType<DataRow>().Skip(1))
            {
                var user = new User
                {
                    Id = Guid.NewGuid().ToString(),
                    UserName = row.ItemArray[0].ToString()?.Trim(),
                    FirstName = row.ItemArray[1].ToString()?.Trim(),
                    LastName = row.ItemArray[2].ToString()?.Trim(),
                    Street = row.ItemArray[3].ToString()?.Trim(),
                    City = row.ItemArray[4].ToString()?.Trim(),
                    Zip = row.ItemArray[5].ToString()?.Trim(),
                    Country = row.ItemArray[6].ToString()?.Trim(),
                    Email = row.ItemArray[7].ToString()?.Trim(),
                    PhoneNumber = row.ItemArray[8].ToString()?.Trim()
                };

                user.Creator = user.Modifier = _authContext.GetCurrentUser()?.UserName;
                user.Created = user.Modified = DateTime.UtcNow;

                var existingPerson =  existingPersons.FirstOrDefault(x =>
                    x.UserName == user.UserName &&
                    x.FirstName == user.FirstName &&
                    x.LastName == user.LastName &&
                    x.Street == user.Street &&
                    x.City == user.City &&
                    x.Zip == user.Zip &&
                    x.Country == user.Country &&
                    x.PhoneNumber == user.PhoneNumber &&
                    x.Email == user.Email);

                if (existingPerson == null)
                {
                    newUsers.Add(user);
                    existingPersons.Add(user);
                }
            }

            await _dbContext.SaveBatch(newUsers);

            return new ImportResultModel
            {
                ImportedPersons = newUsers.Count
            };
        }

        public async Task ClearUsers()
        {
            var result = await _dbContext.GetList<User>();

            foreach (var user in result.List)
            {
                await _dbContext.Delete<User>(user.Id);
            }
        }
    }
}
