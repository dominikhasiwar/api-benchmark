using System.Net;
using DotnetApi.Models;
using Microsoft.AspNetCore.Diagnostics;
using Serilog;

namespace DotnetApi.Handlers
{
    public static class GlobalExceptionHandler
    {
        public static void ConfigureExceptionHandler(this IApplicationBuilder app)
        {
            app.UseExceptionHandler(appError =>
            {
                appError.Run(async context =>
                {
                    context.Response.StatusCode = (int)HttpStatusCode.InternalServerError;
                    context.Response.ContentType = "application/json";

                    var contextFeature = context.Features.Get<IExceptionHandlerFeature>();
                    if (contextFeature != null)
                    {
                        Log.Error(contextFeature.Error.Message, contextFeature.Error);

                        var errorModel = CreateErrorModel(contextFeature.Error);

                        context.Response.StatusCode = (int)errorModel.StatusCode;
                        await context.Response.WriteAsync(errorModel.ToString());
                    }
                });
            });
        }

        private static ErrorModel CreateErrorModel(Exception exception)
        {
            switch (exception)
            {
                default:
                    return new ErrorModel
                    {
                        StatusCode = HttpStatusCode.InternalServerError,
                        ErrorMessage = exception.Message
                    };
            }
        }
    }
}
