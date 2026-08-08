class CreateTracks < ActiveRecord::Migration[8.0]
  def change
    create_table :tracks do |t|
      t.references :composition, null: false, foreign_key: true
      t.string :name, null: false
      t.float :gain, null: false, default: 1.0
      t.float :pan, null: false, default: 0.0
      t.jsonb :midi_data, null: false, default: {}

      t.timestamps
    end
  end
end
