require 'spec_helper'

module ChargeRequest
  describe Qb, :focus do
    let(:transaction_id) { 1 }
    let(:mobile_number) { "85513123456" }

    subject { Qb.new(transaction_id, mobile_number) }

    describe ".save!" do
      it "should do something" do
        VCR.use_cassette(:qb_success) do
          subject.save!
        end
      end
    end
  end
end
