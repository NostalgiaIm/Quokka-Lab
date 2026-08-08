class PlaylistItem < ApplicationRecord
  belongs_to :playlist
  belongs_to :composition

  validates :position, numericality: { greater_than_or_equal_to: 0 }
end
