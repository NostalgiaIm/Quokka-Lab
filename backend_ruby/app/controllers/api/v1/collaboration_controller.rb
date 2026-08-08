module Api
  module V1
    class CollaborationController < ApplicationController
      skip_before_action :authenticate_user!, only: :create

      def create
        event = params.require(:event).permit(:type, :note, :frequency, :at).to_h
        ActionCable.server.broadcast("composition_room_#{params[:room_id]}", { event: event })
        head :accepted
      end
    end
  end
end
