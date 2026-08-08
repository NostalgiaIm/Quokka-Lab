class CreateComments < ActiveRecord::Migration[8.0]
  def change
    create_table :comments do |t|
      t.references :composition, null: false, foreign_key: true
      t.references :user, foreign_key: true
      t.text :body, null: false

      t.timestamps
    end
  end
end
