require 'rails_helper'

describe BeelineChargeRequestUpdaterJob do
  # this job will be queued by another app, therefore only perform method is important
  let(:session_id) { "1234" }
  let(:result_code) { "some-code" }

  let(:options) { {"session_id" => session_id, "result_code" => result_code} }

  subject { described_class.new(options) }

  it "should be serializeable" do
    expect(subject.serialize["arguments"].first).to eq(options)
  end

  describe "#perform(session_id, result_code)" do
    let(:charge_request_result) { double(::ChargeRequestResult::Beeline) }

    before do
      allow(::ChargeRequestResult::Beeline).to receive(:new).with(options.symbolize_keys).and_return(charge_request_result)
      allow(charge_request_result).to receive(:save!)
    end

    it "should try to build and save a new charge request for the operator" do
      expect(charge_request_result).to receive(:save!)
      subject.perform(session_id, result_code)
    end
  end
end
