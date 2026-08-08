module Internal
  class LioraTriggersController < ApplicationController
    skip_before_action :authenticate_user!

    def create
      # 预留给 Liora：可触发音乐事件、AI 和弦建议或虚拟生命反应。
      Rails.logger.info("Liora trigger received: #{params.to_unsafe_h.except(:controller, :action)}")
      render json: { status: "accepted", integration: "liora", received_at: Time.current.iso8601 }, status: :accepted
    end
  end
end
