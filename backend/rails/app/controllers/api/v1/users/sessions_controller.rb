module Api
  module V1
    module Users
      class SessionsController < Devise::SessionsController
        respond_to :json

        private

        def respond_with(resource, _opts = {})
          render json: {
            user: {
              id: resource.id,
              email: resource.email
            },
            message: "Logged in."
          }, status: :ok
        end

        def respond_to_on_destroy
          if current_user
            render json: { message: "Logged out." }, status: :ok
          else
            render json: { message: "No active session." }, status: :unauthorized
          end
        end
      end
    end
  end
end
