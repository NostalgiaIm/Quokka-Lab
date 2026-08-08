Rails.application.configure do
  config.enable_reloading = true
  config.eager_load = false
  config.consider_all_requests_local = true
  config.server_timing = true
  config.active_storage.service = :local
  config.action_cable.url = "ws://localhost:3000/cable"
  config.action_cable.allowed_request_origins = [%r{http://localhost:\d+}]
end
