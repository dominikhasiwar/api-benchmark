using Amazon;
using Amazon.DynamoDBv2;
using Amazon.DynamoDBv2.DataModel;
using Amazon.Lambda.Serialization.SystemTextJson;
using DotnetApi.Context;
using DotnetApi.Extensions;
using DotnetApi.Handlers;
using DotnetApi.Repositories;
using FluentValidation;
using FluentValidation.AspNetCore;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Mvc.Authorization;
using Microsoft.Identity.Web;
using Microsoft.OpenApi.Models;
using Serilog;

namespace DotnetApi
{
    public class Program
    {
        private static readonly IConfiguration Configuration;

        private static readonly IAppConfig AppConfig;

        static Program()
        {
            Configuration = new ConfigurationBuilder()
                .SetBasePath(AppContext.BaseDirectory)
                .AddJsonFile("appsettings.json", optional: true, reloadOnChange: true)
                .AddEnvironmentVariables("BE_")
                .Build();

            AppConfig = Configuration.Get<AppConfig>();
        }

        public static void Main(string[] args)
        {
            Log.Logger = new LoggerConfiguration()
               .Enrich.FromLogContext()
               .ReadFrom.Configuration(Configuration)
               .CreateLogger();

            var builder = WebApplication.CreateBuilder(args);

            builder.Services.AddLogging(cfg =>
            {
                cfg.ClearProviders();
                cfg.AddSerilog(Log.Logger);
            });

            builder.Services.ConfigureHttpJsonOptions(options =>
            {
                options.SerializerOptions.TypeInfoResolver = ApiSerializerContext.Default;
            });

            builder.Services.AddControllers(opt =>
            {
                opt.Filters.Add(new AuthorizeFilter(AppRoles.READERS));
            })
            .AddJsonOptions(opt =>
            {
                opt.JsonSerializerOptions.TypeInfoResolverChain.Add(ApiSerializerContext.Default);
            });

            builder.Services.AddAWSLambdaHosting(LambdaEventSource.HttpApi, options =>
            {
                options.Serializer = new SourceGeneratorLambdaJsonSerializer<ApiSerializerContext>();
            });

            builder.Services.AddEndpointsApiExplorer();

            builder.Services.AddSwaggerGen(c =>
            {
                c.SwaggerDoc("v1", new OpenApiInfo { Title = "Dotnet Demo API", Version = "1.0" });
                c.AddSecurityDefinition("oauth2", new OpenApiSecurityScheme
                {
                    Description = "OAuth2.0 which uses the implicit flow",
                    Name = "oauth2.0",
                    Type = SecuritySchemeType.OAuth2,
                    Flows = new OpenApiOAuthFlows
                    {
                        Implicit = new OpenApiOAuthFlow()
                        {
                            AuthorizationUrl = new Uri($"https://login.microsoftonline.com/{AppConfig.AzureAd.TenantId}/oauth2/v2.0/authorize"),
                            TokenUrl = new Uri($"https://login.microsoftonline.com/{AppConfig.AzureAd.TenantId}/oauth2/v2.0/token"),
                            Scopes = new Dictionary<string, string>()
                            {
                                { AppConfig.AzureAd.Scopes, string.Empty }
                            }
                        }
                    }
                });
                c.AddSecurityRequirement(new OpenApiSecurityRequirement
                {
                    {
                        new OpenApiSecurityScheme
                        {
                            Reference = new OpenApiReference{Type=ReferenceType.SecurityScheme, Id="oauth2"}
                        },
                        new []{ AppConfig.AzureAd.Scopes }
                    }
                  });
            });

            builder.Services.AddSingleton(new DynamoDBContextConfig
            {
                Conversion = DynamoDBEntryConversion.V2
            });

            if (AppConfig.DynamoDb?.ServiceUrl?.HasValue() == true)
            {
                builder.Services.AddSingleton(new AmazonDynamoDBConfig
                {
                    RegionEndpoint = RegionEndpoint.USEast1,
                    ServiceURL = AppConfig.DynamoDb?.ServiceUrl?.HasValue() == true ? AppConfig.DynamoDb?.ServiceUrl : null,
                });
            }
            else
            {
                builder.Services.AddSingleton(new AmazonDynamoDBConfig
                {
                    RegionEndpoint = RegionEndpoint.USEast1,
                });
            }

            builder.Services.AddSingleton<IAmazonDynamoDB, AmazonDynamoDBClient>();

            builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
                .AddMicrosoftIdentityWebApi(Configuration);

            builder.Services.AddAuthorization(opt =>
            {
                opt.AddPolicy(AppRoles.READERS, p => p.RequireRole(AppRoles.READERS, AppRoles.WRITERS, AppRoles.ADMINS));
                opt.AddPolicy(AppRoles.WRITERS, p => p.RequireRole(AppRoles.WRITERS, AppRoles.ADMINS));
                opt.AddPolicy(AppRoles.ADMINS, p => p.RequireRole(AppRoles.ADMINS));
            });

            builder.Services.AddHttpContextAccessor();

            builder.Services.AddSingleton<IAuthContext, AuthContext>();

            builder.Services.AddAutoMapper(typeof(Program).Assembly);

            builder.Services.AddScoped<IDynamoDBContext, DynamoDBContext>();

            builder.Services.AddScoped<IAppDbContext>(s => new AppDbContext(AppConfig.DynamoDb?.TableName, s.GetRequiredService<IDynamoDBContext>(), s.GetRequiredService<IAmazonDynamoDB>()));

            builder.Services.AddScoped<IUserRepository, UserRepository>();

            builder.Services.AddValidatorsFromAssemblyContaining<Program>();

            builder.Services.AddFluentValidationAutoValidation();

            var app = builder.Build();

            app.ConfigureExceptionHandler();

            app.UseAuthentication();

            app.UseAuthorization();

            app.MapControllers();

            app.MapGet("/", () => "Welcome to the dotnet api!").ExcludeFromDescription();

            app.UseSwaggerUI(opt =>
            {
                opt.OAuthClientId(AppConfig.AzureAd.ClientId);
                opt.OAuthUsePkce();
                opt.OAuthScopeSeparator(" ");
                opt.OAuthScopes(AppConfig.AzureAd.Scopes);
            });

            app.UseSwagger();

            app.Run();
        }
    }
}
