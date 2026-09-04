class CreateLikes < ActiveRecord::Migration[8.0]
  def change
    create_table :likes do |t|
      t.references :composition, null: false, foreign_key: true
      t.references :user, foreign_key: true

      t.timestamps
    end

    add_index :likes, %i[user_id composition_id], unique: true
  end
end
