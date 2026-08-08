Rails.application.configure do
  config.enable_reloading = false
  config.eager_load = true
  config.consider_all_requests_local = false
  config.public_file_server.headers = { "Cache-Control" => "public, max-age=31536000" }
  config.active_storage.service = :local
  config.force_ssl = ENV.fetch("FORCE_SSL", "false") == "true"
  config.action_cable.url = ENV.fetch("ACTION_CABLE_URL", "ws://localhost:3000/cable")
end
