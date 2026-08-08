require "rails_helper"

RSpec.describe Composition, type: :model do
  it "requires a title" do
    composition = described_class.new(title: "", bpm: 120, key_signature: "C")

    expect(composition).not_to be_valid
    expect(composition.errors[:title]).to be_present
  end
end
