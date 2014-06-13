require 'spec_helper'

describe ChargeRequestResult::Beeline do
  include OperatorExamples
  include ChargeRequestResultExamples

  ASSERTED_OPERATOR = "beeline"

  let(:asserted_operator) { ASSERTED_OPERATOR }
  let(:initialization_args) { [] }

  it_should_behave_like "an operator"
  it_should_behave_like "a charge request result"

  describe "#save!" do
    include ResqueHelpers

    let(:sample_result) { { :session_id => "foobar", :result_code => "barfoo" } }

    subject { described_class.new(sample_result) }

    before do
      do_background_task(:queue_only => true) { subject.save! }
    end

    it "should queue a job for updating the charge request" do
      assert_chibi_charge_request_updater_job(
        "foobar",
        "result",
        asserted_operator,
        "reason"
      )
    end
  end
end
