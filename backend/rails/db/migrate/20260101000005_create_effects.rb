class CreateEffects < ActiveRecord::Migration[8.0]
  def change
    create_table :effects do |t|
      t.references :track, null: false, foreign_key: true
      t.string :kind, null: false
      t.jsonb :parameters, null: false, default: {}
      t.integer :position, null: false, default: 0

      t.timestamps
    end
  end
end
