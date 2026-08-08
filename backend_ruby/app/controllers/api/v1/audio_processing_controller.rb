module Api
  module V1
    class AudioProcessingController < ApplicationController
      skip_before_action :authenticate_user!, only: :create

      def create
        unless params[:audio].present?
          return render json: { error: "audio file is required" }, status: :bad_request
        end

        input_path = persist_upload(params[:audio])
        output_path = Rails.root.join("tmp/audio/processed-#{SecureRandom.hex(8)}.wav")
        reverb = params.fetch(:reverb, 0.5).to_f.clamp(0.0, 1.0)

        command = [
          engine_binary.to_s,
          input_path.to_s,
          output_path.to_s,
          "--reverb",
          reverb.to_s
        ]

        Rails.logger.info("Running audio engine: #{command.join(' ')}")

        if system(*command)
          composition = Composition.create!(
            title: params.fetch(:title, "Processed audio"),
            bpm: params.fetch(:bpm, 120),
            key_signature: params.fetch(:key_signature, "C"),
            is_public: false,
            user: current_user
          )
          composition.audio.attach(io: File.open(output_path), filename: output_path.basename.to_s, content_type: "audio/wav")

          render json: { audio_url: composition.audio_public_url, composition_id: composition.id }, status: :created
        else
          render json: { error: "audio engine failed" }, status: :unprocessable_entity
        end
      ensure
        FileUtils.rm_f(input_path) if input_path
        FileUtils.rm_f(output_path) if output_path
      end

      private

      def persist_upload(upload)
        FileUtils.mkdir_p(Rails.root.join("tmp/audio"))
        path = Rails.root.join("tmp/audio/input-#{SecureRandom.hex(8)}#{File.extname(upload.original_filename)}")
        File.binwrite(path, upload.read)
        path
      end

      def engine_binary
        ENV.fetch("QUOKKA_AUDIO_ENGINE", Rails.root.join("../audio_engine/build/quokka_audio").to_s)
      end
    end
  end
end
