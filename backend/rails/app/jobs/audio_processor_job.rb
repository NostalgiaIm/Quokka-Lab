class AudioProcessorJob < ApplicationJob
  queue_as :default

  # 把耗时的离线音频处理移出 HTTP 请求周期。
  def perform(composition_id, effect_params = {})
    composition = Composition.find(composition_id)
    Rails.logger.info("Queued audio processing for composition=#{composition.id} params=#{effect_params.inspect}")

    # 后续处理路径：
    # 1. 把原始音频 blob 下载到 tmp。
    # 2. 带着 effect_params 调用 audio_engine。
    # 3. 把处理后的输出作为新 blob 或音轨渲染结果挂回作品。
  end
end
