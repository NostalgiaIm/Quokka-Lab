class Playlist < ApplicationRecord
  belongs_to :user, optional: true
  has_many :playlist_items, dependent: :destroy
  has_many :compositions, through: :playlist_items

  validates :title, presence: true, length: { maximum: 120 }
end
