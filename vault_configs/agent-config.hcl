pid_file = "/tmp/agent.pid"

auto_auth {
  method {
    type = "approle"

    config = {
      role_id_file_path = "/vault/config/role_id"
      secret_id_file_path = "/vault/config/secret_id" 
    }
  }

  sink "file" {
    config = {
      path = "/vault/agent-token"
      perms = "0600"
    }
  }
}

template {
    contents = "DB_USER={{ with secret \"secret/data/dotenv/db\" }}{{ .Data.data.db_user }}{{ end }} DB_PASSWORD={{ with secret \"secret/data/dotenv/db\" }}{{ .Data.data.db_password }}{{ end }} REDIS_URL={{ with secret \"secret/data/dotenv/redis\"}}{{ .Data.data.redis_url }}{{ end }} REDIS_PASS={{ with secret \"secret/data/dotenv/redis\" }}{{ .Data.data.redis_pass }}{{ end }} RABBITMQ_URL={{ with secret \"secret/data/dotenv/rabbitmq\" }}{{ .Data.data.rabbitmq_url }}{{ end }} APP_ENV=production"
    destination = "/vault/app_env/app_env.env"
    perms = "0600"
}