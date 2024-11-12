terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.0"
    }
  }
}

provider "docker" {
  host     = "ssh://{var.ssh_host}:{var.ssh_port}"
  ssh_opts = []
}

locals {
  container_name = var.env == "prod" ? "elrek-system_go_prod" : "elrek-system_go_dev"
  port_mapping   = var.env == "prod" ? "3000:3000" : "3001:3000"
}

resource "docker_image" "elrek_system_go" {
  name = "sc4n1a471/elrek-system_go:${var.container_version}"
}

resource "docker_container" "elrek_system_go" {
  name  = local.container_name
  image = docker_image.elrek_system_go.name

  volumes {
    host_path      = "/var/log/elrek-system_go"
    container_path = "/app/logs"
  }

  env = [
    "DB_USERNAME=${var.db_username}",
    "DB_PASSWORD=${var.db_password}",
    "DB_HOST=${var.db_host}",
    "DB_PORT=${var.db_port}",
    "DB_NAME=${var.db_name}",
    "FRONTEND_URL=${var.frontend_url}",
    "BACKEND_URL=${var.backend_url}",
    "DOMAIN=${var.domain}",
  ]

  ports {
    internal = 3000
    external = var.env == "prod" ? 3000 : 3001
  }

  restart = "on-failure"
  max_retry_count = 5
}