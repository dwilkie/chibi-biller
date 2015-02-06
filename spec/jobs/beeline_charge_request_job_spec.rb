require 'rails_helper'

describe BeelineChargeRequestJob do
  describe "#queue_name" do
    it { expect(subject.queue_name).to eq(Rails.application.secrets[:beeline_charge_request_queue]) }
  end

  it "should be a type of ActiveJob::Base" do
    expect(subject).to be_a(ActiveJob::Base)
  end
end
