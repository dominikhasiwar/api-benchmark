resource "aws_lambda_function" "go_lambda" {
  filename         = "go.zip"
  function_name    = var.aws_go_lambda_function_name
  role             = aws_iam_role.lambda_exec_role.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  memory_size      = 1024
  timeout          = 30
  source_code_hash = filebase64sha256("go.zip")
  environment {
    variables = {
      BE_AWS_REGION        = var.aws_region
      BE_AWS_TABLE_NAME    = var.aws_dynamodb_table_name
      BE_AUTH_TENANT_ID    = var.azure_tenant_id
      BE_AUTH_CLIENT_ID    = var.azure_client_id
      BE_AUTH_SCOPE        = var.azure_auth_scope
      BE_SWAGGER_BASE_PATH = "/default"
    }
  }
}

resource "aws_lambda_function" "dotnet_lambda" {
  filename         = "dotnet.zip"
  function_name    = var.aws_dotnet_lambda_function_name
  role             = aws_iam_role.lambda_exec_role.arn
  handler          = "DotnetApi"
  runtime          = "dotnet8"
  memory_size      = 1024
  timeout          = 30
  source_code_hash = filebase64sha256("dotnet.zip")
}
