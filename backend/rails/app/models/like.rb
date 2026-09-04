class Like < ApplicationRecord
  belongs_to :composition
  belongs_to :user, optional: true

  validates :user_id, uniqueness: { scope: :composition_id }, allow_nil: true
end
