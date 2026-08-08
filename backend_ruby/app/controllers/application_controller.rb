class ApplicationController < ActionController::API
  include ActionController::MimeResponds

  before_action :authenticate_user!, unless: :public_endpoint?

  private

  # MVP 阶段允许未登录用户浏览和预览公开作品。
  def public_endpoint?
    controller_path == "api/v1/compositions" && action_name.in?(%w[index show create])
  end
end
