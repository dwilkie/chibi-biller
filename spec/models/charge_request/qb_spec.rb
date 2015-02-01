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
    include WebMockHelpers
    include ChargeRequestHelpers
    include ActiveJobHelpers

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
      expect(last_request(:url)).to eq(asserted_url)
      expect(last_request.method).to eq(:get)
    end

    context "qb responds with 200 OK" do
      before do
        expect_charge_request { subject.save! }
      end

      it "should do nothing" do
        expect(chibi_charge_request_updater_job).to be_nil
      end
    end

    context "qb responds with 500 Internal Server Error" do
      before do
        trigger_job(:queue_only => true) do
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
