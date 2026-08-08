module Api
  module V1
    class PlaylistsController < ApplicationController
      skip_before_action :authenticate_user!, only: %i[index show]

      def index
        render json: Playlist.where(is_public: true).order(created_at: :desc).as_json(include: :compositions)
      end

      def show
        render json: Playlist.find(params[:id]).as_json(include: :compositions)
      end

      def create
        playlist = Playlist.new(playlist_params)
        playlist.user = current_user if user_signed_in?

        if playlist.save
          render json: playlist, status: :created
        else
          render json: { errors: playlist.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def update
        playlist = Playlist.find(params[:id])

        if playlist.update(playlist_params)
          render json: playlist
        else
          render json: { errors: playlist.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def destroy
        Playlist.find(params[:id]).destroy
        head :no_content
      end

      private

      def playlist_params
        params.require(:playlist).permit(:title, :description, :is_public)
      end
    end
  end
end
