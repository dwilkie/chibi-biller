require 'rails_helper'

describe ChargeRequest::Beeline do
  include OperatorExamples
  include ChargeRequestExamples

  let(:transaction_id) { 2 }
  let(:mobile_number) { "85560201158" }
  let(:asserted_operator) { "beeline" }
  let(:initialization_args) { [transaction_id, mobile_number] }

  subject { described_class.new(*initialization_args) }

  it_should_behave_like "an operator"
  it_should_behave_like "a charge request"

  describe "#save!" do
    include ResqueHelpers

    before do
      subject.save!
    end

    it "should schedule a job to send the charge request to Beeline" do
      queue = ResqueSpec.queues[ENV["BEELINE_CHARGE_REQUEST_QUEUE"]].first
      expect(queue).not_to be_nil
      expect(queue[:class]).to eq(ENV["BEELINE_CHARGE_REQUEST_WORKER"])
      expect(queue[:args]).to eq([transaction_id.to_s, mobile_number])
    end
  end
end
