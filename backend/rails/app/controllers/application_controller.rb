class ApplicationController < ActionController::API
  include ActionController::MimeResponds

  before_action :authenticate_user!, unless: :public_endpoint?

  private

  # MVP mode allows guests to browse, preview, and create public compositions.
  def public_endpoint?
    controller_path == "api/v1/compositions" && action_name.in?(%w[index show create])
  end
end
