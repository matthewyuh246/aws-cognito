variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-1"
}

variable "domain" {
  description = "Domain name for the application"
  type        = string
  default     = "localhost"
}

variable "cognito_domain" {
  description = "Cognito domain prefix"
  type        = string
  default     = "social-login-app"
}

variable "google_client_id" {
  description = "Google OAuth Client ID"
  type        = string
  sensitive   = true
}

variable "google_client_secret" {
  description = "Google OAuth Client Secret"
  type        = string
  sensitive   = true
}