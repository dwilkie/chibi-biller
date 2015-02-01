require 'rails_helper'

describe BeelineChargeRequestJob do
  # this job will be run by another app therefore only the queue name is important
  describe "#queue_name" do
    it { expect(subject.queue_name).to eq("beeline_charge_request_queue") }
  end

  it "should be a type of ActiveJob::Base" do
    expect(subject).to be_a(ActiveJob::Base)
  end
end
