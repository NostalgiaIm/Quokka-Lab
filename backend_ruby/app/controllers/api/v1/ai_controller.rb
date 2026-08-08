module Api
  module V1
    class AiController < ApplicationController
      skip_before_action :authenticate_user!, only: %i[chords mix style_transfer]

      def chords
        # 外部 AI 模型接入点：把旋律事件映射为和弦建议。
        render json: {
          suggestions: [
            { chord: "Cmaj7", confidence: 0.72 },
            { chord: "Am7", confidence: 0.58 },
            { chord: "Fmaj7", confidence: 0.51 }
          ]
        }
      end

      def mix
        # 智能混音接入点：给出避免削波的增益建议。
        render json: {
          master_gain: 0.86,
          track_gains: params.fetch(:tracks, []).map.with_index { |_track, index| { index: index, gain: 0.9 } }
        }
      end

      def style_transfer
        # 接入真实风格迁移模型后，这里应改为异步任务。
        render json: { status: "queued", style: params.fetch(:style, "jazz") }, status: :accepted
      end
    end
  end
end
