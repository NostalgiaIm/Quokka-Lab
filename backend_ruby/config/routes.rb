Rails.application.routes.draw do
  mount ActionCable.server => "/cable"

  devise_for :users, defaults: { format: :json }, skip: %i[sessions registrations]

  devise_scope :user do
    post "api/v1/signup", to: "api/v1/users/registrations#create"
    post "api/v1/login", to: "api/v1/users/sessions#create"
    delete "api/v1/logout", to: "api/v1/users/sessions#destroy"
  end

  namespace :api do
    namespace :v1 do
      resources :compositions do
        resources :comments, only: %i[index create destroy]
        resource :like, only: %i[create destroy]
      end
      resources :playlists
      post "ai/chords", to: "ai#chords"
      post "ai/mix", to: "ai#mix"
      post "ai/style_transfer", to: "ai#style_transfer"
      post "audio/process", to: "audio_processing#create"
      post "collaboration/:room_id/events", to: "collaboration#create"
    end
  end

  namespace :internal do
    post "liora_trigger", to: "liora_triggers#create"
  end
end
