Devise.setup do |config|
  config.mailer_sender = "quokka@example.com"
  require "devise/orm/active_record"

  config.jwt do |jwt|
    jwt.secret = ENV.fetch("DEVISE_JWT_SECRET_KEY", "change-me-in-development")
    jwt.dispatch_requests = [["POST", %r{^/api/v1/login$}]]
    jwt.revocation_requests = [["DELETE", %r{^/api/v1/logout$}]]
    jwt.expiration_time = 1.day.to_i
  end
end
