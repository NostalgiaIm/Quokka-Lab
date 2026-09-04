class CreateCompositions < ActiveRecord::Migration[8.0]
  def change
    create_table :compositions do |t|
      t.string :title, null: false
      t.references :user, foreign_key: true
      t.string :audio_url
      t.jsonb :midi_data, null: false, default: {}
      t.integer :bpm, null: false, default: 120
      t.string :key_signature, null: false, default: "C"
      t.boolean :is_public, null: false, default: true

      t.timestamps
    end

    add_index :compositions, :is_public
    add_index :compositions, :midi_data, using: :gin
  end
end
