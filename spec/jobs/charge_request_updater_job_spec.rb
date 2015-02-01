require 'rails_helper'

describe ChargeRequestUpdaterJob do
  # this job will be run on Chibi therefore only the queue name is important
  describe "#queue_name" do
    it { expect(subject.queue_name).to eq("critical") }
  end

  it "should be a type of ActiveJob::Base" do
    expect(subject).to be_a(ActiveJob::Base)
  end
end
