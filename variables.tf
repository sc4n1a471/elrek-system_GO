variable "container_version" {
  description = "Version of the Docker image"
  type        = string
  default     = "latest-dev"
}

variable "env" {
  description = "Environment to deploy (prod or dev)"
  type        = string
  default     = "dev"
}

variable "db_username" {
  description = "The database username"
  type        = string
  sensitive = true
}

variable "db_password" {
  description = "The database password"
  type        = string
  sensitive   = true
}

variable "db_host" {
  description = "The database host"
  type        = string
  sensitive = true
}

variable "db_port" {
  description = "The database port"
  type        = string
}

variable "db_name" {
  description = "The database name"
  type        = string
}

variable "domain" {
  description = "The domain name"
  type        = string
  sensitive = true
}

variable "frontend_url" {
  description = "The frontend URL"
  type        = string
  sensitive = true
}

variable "backend_url" {
  description = "The backend URL"
  type        = string
  sensitive = true
}

variable "ssh_host" {
  description = "The SSH host in a <username>@<hostname> format"
  type        = string
  sensitive = true
}

variable "graylog_host" {
  description = "The Graylog host with ip and port: <ip>:<port>"
  type        = string
}