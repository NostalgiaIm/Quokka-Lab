module Api
  module V1
    class CompositionsController < ApplicationController
      before_action :set_composition, only: %i[show update destroy]

      def index
        compositions = Composition.publicly_visible.includes(audio_attachment: :blob)
        render json: compositions.map { |composition| serialize(composition) }
      end

      def show
        render json: serialize(@composition)
      end

      def create
        composition = Composition.new(composition_params.except(:audio))
        composition.user = current_user if user_signed_in?
        attach_audio(composition, composition_params[:audio])

        if composition.save
          render json: serialize(composition), status: :created
        else
          render json: { errors: composition.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def update
        if @composition.update(composition_params.except(:audio))
          attach_audio(@composition, composition_params[:audio])
          @composition.save if composition_params[:audio].present?
          render json: serialize(@composition)
        else
          render json: { errors: @composition.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def destroy
        @composition.destroy
        head :no_content
      end

      private

      def set_composition
        @composition = Composition.find(params[:id])
      end

      def composition_params
        permitted = params.require(:composition).permit(
          :title,
          :audio,
          :audio_url,
          :bpm,
          :key_signature,
          :is_public,
          :midi_data,
          midi_data: {}
        )
        permitted[:midi_data] = parse_midi_data(permitted[:midi_data])
        permitted
      end

      def parse_midi_data(value)
        return {} if value.blank?
        return value if value.is_a?(Hash)

        JSON.parse(value)
      rescue JSON::ParserError
        {}
      end

      def attach_audio(composition, upload)
        return unless upload.present?

        composition.audio.attach(upload)
      end

      def serialize(composition)
        {
          id: composition.id,
          title: composition.title,
          audio_url: composition.audio_public_url,
          midi_data: composition.midi_data || {},
          bpm: composition.bpm,
          key_signature: composition.key_signature,
          is_public: composition.is_public,
          likes_count: composition.likes.count,
          comments_count: composition.comments.count,
          created_at: composition.created_at,
          updated_at: composition.updated_at
        }
      end
    end
  end
end
