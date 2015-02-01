require 'rails_helper'

describe ChargeRequest::Qb do
  include OperatorExamples
  include ChargeRequestExamples

  let(:transaction_id) { "2" }
  let(:mobile_number) { "85513123456" }
  let(:asserted_operator) { "qb" }
  let(:initialization_args) { [transaction_id, mobile_number] }

  subject { described_class.new(*initialization_args) }

  it_should_behave_like "an operator"
  it_should_behave_like "a charge request"

  describe "#save!" do
    include ResqueHelpers
    include WebMockHelpers
    include ChargeRequestHelpers

    ASSERTED_TRANSACTION_ID_KEY = "TRANID"
    ASSERTED_MOBILE_NUMBER_KEY = "MSISDN"

    def asserted_url
      @asserted_url ||= super(
        asserted_operator, :query => {
          ASSERTED_MOBILE_NUMBER_KEY => mobile_number,
          ASSERTED_TRANSACTION_ID_KEY => transaction_id
        }
      )
    end

    def expect_charge_request(options = {}, &block)
      super({:operator => asserted_operator, :url => asserted_url}.merge(options), &block)
    end

    it "should send the charge request to qb" do
      expect_charge_request { subject.save! }
      last_request(:url).should == asserted_url
      last_request.method.should == :get
    end

    context "qb responds with 200 OK" do
      before do
        do_background_task(:queue_only => true) do
          expect_charge_request { subject.save! }
        end
      end

      it "should do nothing" do
        chibi_charge_request_updater_job.should be_nil
      end
    end

    context "qb responds with 500 Internal Server Error" do
      before do
        do_background_task(:queue_only => true) do
          expect_charge_request(:response_type => :internal_server_error) { subject.save! }
        end
      end

      it "should queue a job for updating the charge request with an 'errored' result" do
        assert_chibi_charge_request_updater_job(
          transaction_id, "errored", asserted_operator, "Client Error: Internal Server Error"
        )
      end
    end
  end
end
