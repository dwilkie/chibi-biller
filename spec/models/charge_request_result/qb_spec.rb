require 'spec_helper'

module ChargeRequestResult
  describe Qb do
    include OperatorExamples
    include ChargeRequestResultExamples

    ASSERTED_OPERATOR = "qb"
    let(:asserted_operator) { ASSERTED_OPERATOR }
    let(:initialization_args) { [] }

    it_should_behave_like "an operator"
    it_should_behave_like "a charge request result"

    describe "#save!" do
      include ResqueHelpers

      ASSERTED_ID_KEY = "TRANID"
      ASSERTED_RESULT_KEY = "RESULT"
      ASSERTED_REASON_KEY = "REASON"

      SAMPLE_ID = "1"

      SAMPLE_RESULT_PARAMS = [
        {:result => {:actual => "Successful.", :expected => "successful"}},
        {:result => {:actual => "Failed.", :expected => "failed"}},
        {:result => {:actual => "Error.", :expected => "errored"}}
      ]

      def sample_result_params(sample_result)
        actual_result_params = {
          ASSERTED_ID_KEY => SAMPLE_ID
        }.merge(
          ASSERTED_RESULT_KEY => sample_result[:result][:actual]
        )
        sample_result[:reason] ? actual_result_params.merge(
          ASSERTED_REASON_KEY => sample_result[:reason][:actual]
        ) : actual_result_params
      end

      SAMPLE_RESULT_PARAMS.each do |sample_result|
        actual_reason_string = ", '#{ASSERTED_REASON_KEY}' => '#{sample_result[:reason][:actual]}'" if sample_result[:reason]
        context_string = "{'#{ASSERTED_ID_KEY}' => '#{SAMPLE_ID}', '#{ASSERTED_RESULT_KEY}' => '#{sample_result[:result][:actual]}'#{actual_reason_string}}"
        context "#params => #{context_string}" do
          subject { Qb.new(sample_result_params(sample_result)) }

          before do
            do_background_task(:queue_only => true) { subject.save! }
          end

          expected_reason_string = ", #{(sample_result[:reason] || {})[:expected]}" if sample_result[:reason]

          it "should queue a job for updating the charge request with args: '#{SAMPLE_ID}', '#{sample_result[:result][:expected]}', '#{ASSERTED_OPERATOR}'#{expected_reason_string}" do
            assert_chibi_charge_request_updater_job(
              SAMPLE_ID,
              sample_result[:result][:expected],
              asserted_operator,
              (sample_result[:reason] || {})[:expected]
            )
          end
        end
      end
    end
  end
end
