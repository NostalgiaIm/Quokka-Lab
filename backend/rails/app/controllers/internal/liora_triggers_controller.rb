module Internal
  class LioraTriggersController < ApplicationController
    skip_before_action :authenticate_user!

    def create
      # Reserved for Liora sound-module events, AI chord suggestions, or avatar reactions.
      Rails.logger.info("Liora trigger received: #{params.to_unsafe_h.except(:controller, :action)}")
      render json: { status: "accepted", integration: "liora", received_at: Time.current.iso8601 }, status: :accepted
    end
  end
end
