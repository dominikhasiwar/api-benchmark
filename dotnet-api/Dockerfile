# Official microsoft .NET SDK image
FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS api-build

WORKDIR /var/task

COPY . .

ENV \
    # Enable detection of running in a container
    DOTNET_RUNNING_IN_CONTAINER=true \
    # Lambda is optionated about installing tooling under /var
    DOTNET_ROOT=/var/lang/dotnet \
    # Don't display welcome message on first run
    DOTNET_NOLOGO=true \
    # Disable Microsoft's telemetry collection
    DOTNET_CLI_TELEMETRY_OPTOUT=true


RUN dotnet publish -c Release -r linux-musl-x64 \
   /p:AWSProjectType="Lambda" \
   /p:LangVersion="latest" \
   /p:TargetFramework="net8.0" \
   /p:AssemblyName="bootstrap" \
   /p:RuntimeIdentifier="linux-musl-x64" \
   /p:PublishSingleFile="true" \
   /p:PublishReadyToRun="true" \
   /p:GenerateDocumentationFile="false" \
   /p:ExecutableOutputType="true" \
   /p:OutputType="Exe" \
   /p:CopyLocalLockFileAssemblies="true" \
   /p:PublishTrimmed="true" \
   /p:Optimize="true" \
   /p:TrimMode="link" \
   /p:TieredCompilationQuickJit="false" \
   /p:TieredCompilation="false" \
   /p:CopyLocalLockFileAssemblies="true" \
   /p:SuppressTrimAnalysisWarnings="true" \
  --self-contained true \
  -o release

# Set very small linux distribution as an base for a service
FROM alpine:latest AS runtime

# Set environment from arguments or let default
ARG DEFAULT_ENV
ENV STAGE_ENVIRONMENT=$DEFAULT_ENV

ENV \
    # Enable detection of running in a container
    DOTNET_RUNNING_IN_CONTAINER=true \
    # Lambda is opinionated about installing tooling under /var
    DOTNET_ROOT=/var/lang/bin \
    # Don't display welcome message on first run
    DOTNET_NOLOGO=true \
    # Disable Microsoft's telemetry collection
    DOTNET_CLI_TELEMETRY_OPTOUT=true \
    # The AWS base images provide the following environment variables:
    LAMBDA_TASK_ROOT=/var/task \
    LAMBDA_RUNTIME_DIR=/var/runtime \
    # https://docs.aws.amazon.com/lambda/latest/dg/configuration-concurrency.html?icmpid=docs_lambda_console
    # For the .NET 3.1 runtime, set this variable to enable or disable .NET 3.1 specific runtime optimizations.
    # Values include "always", "never", and "provisioned-concurrency".
    # For information, see Configuring provisioned concurrency.
    AWS_LAMBDA_DOTNET_PREJIT="Always"
    # https://docs.aws.amazon.com/lambda/latest/dg/configuration-concurrency.html?icmpid=docs_lambda_console
    #AWS_LAMBDA_INITIALIZATION_TYPE="provisioned-concurrency"

# Allow to redirect and get work .NET web service from any host
# But needs to open(dispose) ports by run docker
#ENV ASPNETCORE_URLS=http://+:8080

# Expose HTTP port
#EXPOSE 8080

#ENV PATH=/var/lang/bin:/usr/local/bin:/usr/bin/:/bin:/opt/bin
#ENV LD_LIBRARY_PATH=/var/lang/lib:/lib64:/usr/lib64:/var/runtime:/var/runtime/lib:/var/task:/var/task/lib:/opt/lib
#ENV LAMBDA_TASK_ROOT=/var/task
#ENV LAMBDA_RUNTIME_DIR=/var/runtime

# Install dependencies
# https://docs.microsoft.com/en-us/dotnet/core/install/linux-alpine
# --no-cache option allows to not cache the index locally, which is useful for keeping containers small
# Literally it equals `apk update` in the beginning and `rm -rf /var/cache/apk/*` in the end.
RUN apk add --no-cache musl icu-libs krb5-libs

# Set directory to run from
WORKDIR /var/task

# Copy executable
COPY --from=api-build /var/task/release/bootstrap* ./
COPY --from=api-build /var/task/release/appsettings.json ./
#COPY --from=api-build /var/task/bin/entry.sh ./
#COPY --from=api-build /var/task/bin/aws-lambda-rie /usr/bin/aws-lambda-rie
#COPY --from=api-build /var/task/bin/entry.sh /var/task

# (Optional) Add Lambda Runtime Interface Emulator and use a script in the ENTRYPOINT for simpler local runs
#ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
#COPY entry.sh /
#RUN chmod 755 /usr/bin/aws-lambda-rie /entry.sh

#https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-runtime
#aws-proxy(bootsrap)

EXPOSE 5000

ENTRYPOINT ["/var/task/bootstrap"]