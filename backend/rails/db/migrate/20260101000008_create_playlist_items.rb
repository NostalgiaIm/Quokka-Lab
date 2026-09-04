class CreatePlaylistItems < ActiveRecord::Migration[8.0]
  def change
    create_table :playlist_items do |t|
      t.references :playlist, null: false, foreign_key: true
      t.references :composition, null: false, foreign_key: true
      t.integer :position, null: false, default: 0

      t.timestamps
    end

    add_index :playlist_items, %i[playlist_id composition_id], unique: true
  end
end
