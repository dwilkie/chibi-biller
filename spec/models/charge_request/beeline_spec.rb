require 'rails_helper'

describe ChargeRequest::Beeline do
  include OperatorExamples
  include ChargeRequestExamples
  include ActiveJobHelpers

  let(:transaction_id) { 2 }
  let(:mobile_number) { "85560201158" }
  let(:asserted_operator) { "beeline" }
  let(:initialization_args) { [transaction_id, mobile_number] }

  subject { described_class.new(*initialization_args) }

  it_should_behave_like "an operator"
  it_should_behave_like "a charge request"

  describe "#save!" do
    before do
      trigger_job(:queue_only => true) { subject.save! }
    end

    it "should schedule a job to send the charge request to Beeline" do
      job = enqueued_jobs.first
      expect(job).not_to be_nil
      expect(job[:args]).to eq(
        [
          transaction_id.to_s,
          mobile_number,
          BeelineChargeRequestUpdaterJob.queue_name,
          BeelineChargeRequestUpdaterJob.to_s,
          Rails.application.secrets[:beeline_billing_server_address]
        ]
      )
    end
  end
end
