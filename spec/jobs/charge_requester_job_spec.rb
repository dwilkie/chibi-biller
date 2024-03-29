require 'rails_helper'

describe ChargeRequesterJob do
  let(:asserted_queue) { :charge_requester_queue }
  let(:transaction_id) { 1 }
  let(:mobile_number) { "85513123456" }
  let(:operator) { "qb" }

  let(:args) { [transaction_id, operator, mobile_number] }

  subject { described_class.new(*args) }

  it "should be serializeable" do
    expect(subject.serialize["arguments"]).to eq(args)
  end

  describe "#perform(charge_request_id, operator, mobile_number)" do
    let(:charge_request) { double(ChargeRequest::Qb) }

    before do
      allow(ChargeRequest::Qb).to receive(:new).and_return(charge_request)
      allow(charge_request).to receive(:save!)
    end

    it "should try to build and save a new charge request for the operator" do
      expect(ChargeRequest::Qb).to receive(:new).with(transaction_id, mobile_number)
      expect(charge_request).to receive(:save!)
      subject.perform(transaction_id, operator, mobile_number)
    end
  end

  describe "performing the job" do
    # this is an integration test

    include WebMockHelpers
    include ChargeRequestHelpers
    include ActiveJobHelpers

    before do
      ActiveJob::Base.queue_adapter = :test
    end

    def enqueue_job
      trigger_job { described_class.perform_later(transaction_id, operator, mobile_number) }
    end

    def asserted_url
      @asserted_url ||= super(
        operator, :query => {
          "MSISDN" => mobile_number,
          "TRANID" => transaction_id
        }
      )
    end

    it "should send the charge request" do
      expect_charge_request(:operator => operator, :url => asserted_url) { enqueue_job }
      expect(last_request(:method)).to eq(:get)
      expect(last_request(:url)).to eq(asserted_url)
    end
  end
end
