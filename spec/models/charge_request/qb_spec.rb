require 'spec_helper'

module ChargeRequest
  describe Qb, :focus do
    let(:transaction_id) { 1 }
    let(:mobile_number) { "85513123456" }

    ASSERTED_TRANSACTION_ID_KEY = "TRANID"
    ASSERTED_MOBILE_NUMBER_KEY = "MSISDN"

    subject { Qb.new(transaction_id, mobile_number) }

    describe ".save!" do
      include WebMockHelpers

      def asserted_url
        url = URI.parse(ENV["CHARGE_REQUEST_URL_QB"])
        url.query = Rack::Utils.build_query(
          ASSERTED_MOBILE_NUMBER_KEY => mobile_number,
          ASSERTED_TRANSACTION_ID_KEY => transaction_id
        )
        url.to_s
      end

      it "should send the billing request to qb" do
        VCR.use_cassette(:qb_success, :erb => {:url => asserted_url}) do
          subject.save!
          last_request(:url).should == asserted_url
          last_request.method.should == :get
        end
      end
    end
  end
end
