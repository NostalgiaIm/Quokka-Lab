class Effect < ApplicationRecord
  belongs_to :track

  validates :kind, presence: true, inclusion: { in: %w[reverb delay pitch_shift eq compressor] }
  validates :position, numericality: { greater_than_or_equal_to: 0 }
end
