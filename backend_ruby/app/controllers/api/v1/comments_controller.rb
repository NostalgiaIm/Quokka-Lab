module Api
  module V1
    class CommentsController < ApplicationController
      skip_before_action :authenticate_user!, only: %i[index create]

      def index
        composition = Composition.find(params[:composition_id])
        render json: composition.comments.order(created_at: :desc)
      end

      def create
        composition = Composition.find(params[:composition_id])
        comment = composition.comments.new(comment_params)
        comment.user = current_user if user_signed_in?

        if comment.save
          render json: comment, status: :created
        else
          render json: { errors: comment.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def destroy
        Comment.find(params[:id]).destroy
        head :no_content
      end

      private

      def comment_params
        params.require(:comment).permit(:body)
      end
    end
  end
end
