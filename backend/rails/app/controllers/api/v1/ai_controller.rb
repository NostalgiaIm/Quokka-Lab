module Api
  module V1
    class AiController < ApplicationController
      skip_before_action :authenticate_user!, only: %i[chords mix style_transfer]

      def chords
        # Internal AI extension point for mapping melody events to chord suggestions.
        render json: {
          suggestions: [
            { chord: "Cmaj7", confidence: 0.72 },
            { chord: "Am7", confidence: 0.58 },
            { chord: "Fmaj7", confidence: 0.51 }
          ]
        }
      end

      def mix
        # Internal smart-mix extension point for headroom-aware gain suggestions.
        render json: {
          master_gain: 0.86,
          track_gains: params.fetch(:tracks, []).map.with_index { |_track, index| { index: index, gain: 0.9 } }
        }
      end

      def style_transfer
        # This should become an async job after a real style-transfer model is connected.
        render json: { status: "queued", style: params.fetch(:style, "jazz") }, status: :accepted
      end
    end
  end
end
