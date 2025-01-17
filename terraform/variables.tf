variable "aws_region" {
  description = "The AWS region"
  type        = string
  default     = "us-east-1"
}

variable "aws_go_lambda_function_name" {
  description = "The AWS GO Lambda Function Name"
  type        = string
  default     = "go-api"
}

variable "aws_dotnet_lambda_function_name" {
  description = "The AWS Dotnet Lambda Function Name"
  type        = string
  default     = "dotnet-api"
}

variable "aws_dynamodb_table_name" {
  description = "The AWS DynamoDb Table Name"
  type        = string
}

variable "azure_tenant_id" {
  description = "The Azure Tenant ID"
  type        = string
}

variable "azure_client_id" {
  description = "The Azure Client ID"
  type        = string
}

variable "azure_auth_scope" {
  description = "The Azure Auth Scope"
  type        = string
}
