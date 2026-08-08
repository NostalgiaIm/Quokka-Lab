class Track < ApplicationRecord
  belongs_to :composition
  has_many :effects, dependent: :destroy

  validates :name, presence: true
  validates :gain, numericality: { greater_than_or_equal_to: 0, less_than_or_equal_to: 2 }
  validates :pan, numericality: { greater_than_or_equal_to: -1, less_than_or_equal_to: 1 }
end
