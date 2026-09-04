module Api
  module V1
    class LikesController < ApplicationController
      skip_before_action :authenticate_user!, only: %i[create destroy]

      def create
        composition = Composition.find(params[:composition_id])
        like = composition.likes.find_or_initialize_by(user: current_user)

        if like.save
          render json: { likes_count: composition.likes.count }, status: :created
        else
          render json: { errors: like.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def destroy
        composition = Composition.find(params[:composition_id])
        composition.likes.where(user: current_user).destroy_all
        render json: { likes_count: composition.likes.count }
      end
    end
  end
end
