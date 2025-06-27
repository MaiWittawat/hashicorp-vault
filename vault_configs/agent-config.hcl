pid_file = "/tmp/vault-agent.pid"

auto_auth {
  method "approle" {
    mount_path = "auth/approle"
    config = {
      role_id_file_path   = "/vault/config/role_id"
      secret_id_file_path = "/vault/config/secret_id"
    }
  }

  sink "file" {
    config = {
      path = "/tmp/vault-token"
      mode = 0600
    }
  }
}

cache {
  use_auto_auth_token = true
}

listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = true
}
