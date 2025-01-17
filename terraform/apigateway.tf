// ==================================
// Api Gateway for dotnet lambda
// ==================================
resource "aws_apigatewayv2_api" "dotnet_http_api" {
  name          = "dotnet_http_api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "dotnet_default_stage" {
  api_id      = aws_apigatewayv2_api.dotnet_http_api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_apigatewayv2_integration" "dotnet_lambda_integration" {
  api_id                 = aws_apigatewayv2_api.dotnet_http_api.id
  integration_type       = "AWS_PROXY"
  connection_type        = "INTERNET"
  description            = "Lambda example"
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.dotnet_lambda.invoke_arn
  passthrough_behavior   = "WHEN_NO_MATCH"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "dotnet_lambda_proxy_route" {
  api_id    = aws_apigatewayv2_api.dotnet_http_api.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.dotnet_lambda_integration.id}"
}

resource "aws_lambda_permission" "dotnet_api_gateway_permission" {
  statement_id  = "AllowExecutionFromApiGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.dotnet_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.dotnet_http_api.execution_arn}/*/*"
}

// ==================================
// Api Gateway for go lambda
// ==================================
resource "aws_api_gateway_rest_api" "go_rest_api" {
  name               = "go_rest_api"
  description        = "API Gateway for Go Lambda function"
  binary_media_types = ["multipart/form-data"]
}

resource "aws_api_gateway_method" "go_root_method" {
  rest_api_id   = aws_api_gateway_rest_api.go_rest_api.id
  resource_id   = aws_api_gateway_rest_api.go_rest_api.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "go_root_integration" {
  rest_api_id             = aws_api_gateway_rest_api.go_rest_api.id
  resource_id             = aws_api_gateway_rest_api.go_rest_api.root_resource_id
  http_method             = aws_api_gateway_method.go_root_method.http_method
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = aws_lambda_function.go_lambda.invoke_arn
}

resource "aws_api_gateway_resource" "go_proxy_resource" {
  rest_api_id = aws_api_gateway_rest_api.go_rest_api.id
  parent_id   = aws_api_gateway_rest_api.go_rest_api.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "go_proxy_method" {
  rest_api_id   = aws_api_gateway_rest_api.go_rest_api.id
  resource_id   = aws_api_gateway_resource.go_proxy_resource.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "go_proxy_integration" {
  rest_api_id             = aws_api_gateway_rest_api.go_rest_api.id
  resource_id             = aws_api_gateway_resource.go_proxy_resource.id
  http_method             = aws_api_gateway_method.go_proxy_method.http_method
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = aws_lambda_function.go_lambda.invoke_arn
}

resource "aws_api_gateway_deployment" "go_api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.go_rest_api.id
  depends_on = [
    aws_api_gateway_integration.go_root_integration,
    aws_api_gateway_integration.go_proxy_integration
  ]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "go_api_stage" {
  rest_api_id   = aws_api_gateway_rest_api.go_rest_api.id
  deployment_id = aws_api_gateway_deployment.go_api_deployment.id
  stage_name    = "default"
}

resource "aws_lambda_permission" "go_api_gateway_permission" {
  statement_id  = "AllowExecutionFromApiGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.go_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.go_rest_api.execution_arn}/*/*/"
}
