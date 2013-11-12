require 'spec_helper'

describe ChargeRequester do
  let(:asserted_queue) { :charge_requester_queue }
  let(:transaction_id) { 1 }
  let(:mobile_number) { "85513123456" }
  let(:operator) { :qb }

  describe "@queue" do
    it "should == :charge_requester_queue" do
      subject.class.instance_variable_get(:@queue).should == asserted_queue
    end
  end

  describe ".perform(charge_request_id, operator, mobile_number)" do
    let(:charge_request) { double(ChargeRequest::Qb) }

    before do
      ChargeRequest::Qb.stub(:new).and_return(charge_request)
      charge_request.stub(:save!)
    end

    it "should try to build and save a new charge request for the operator" do
      ChargeRequest::Qb.should_receive(:new).with(transaction_id, mobile_number)
      charge_request.should_receive(:save!)
      subject.class.perform(transaction_id, operator, mobile_number)
    end
  end

  describe "performing the job" do
    # this is an integration test

    include ResqueHelpers
    include WebMockHelpers
    include ChargeRequestHelpers

    def enqueue_job
      Resque.enqueue(ChargeRequester, transaction_id, operator, mobile_number)
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
      do_background_task(:queue_only => true) { enqueue_job }
      expect_charge_request(:operator => operator, :url => asserted_url) { perform_background_job(asserted_queue) }
      last_request(:method).should == :get
      last_request(:url).should == asserted_url
    end
  end
end
