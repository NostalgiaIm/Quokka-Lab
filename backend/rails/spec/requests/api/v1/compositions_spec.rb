require "rails_helper"

RSpec.describe "Compositions", type: :request do
  describe "GET /api/v1/compositions" do
    it "returns public compositions" do
      Composition.create!(title: "First take", bpm: 120, key_signature: "C", is_public: true)

      get "/api/v1/compositions"

      expect(response).to have_http_status(:ok)
      expect(JSON.parse(response.body).first["title"]).to eq("First take")
    end
  end
end
