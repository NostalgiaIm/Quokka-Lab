class Composition < ApplicationRecord
  belongs_to :user, optional: true

  has_many :tracks, dependent: :destroy
  has_many :effects, through: :tracks
  has_many :comments, dependent: :destroy
  has_many :likes, dependent: :destroy

  has_one_attached :audio

  validates :title, presence: true, length: { maximum: 120 }
  validates :bpm, numericality: { greater_than: 20, less_than_or_equal_to: 300 }
  validates :key_signature, presence: true

  scope :publicly_visible, -> { where(is_public: true).order(created_at: :desc) }

  # Returns a browser-playable audio URL; production can switch to S3 or MinIO.
  def audio_public_url
    audio.attached? ? Rails.application.routes.url_helpers.rails_blob_path(audio, only_path: true) : audio_url
  end
end
