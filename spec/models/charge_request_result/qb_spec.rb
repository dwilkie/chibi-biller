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

      ASSERTED_TRANSACTION_ID_PARAM = "TRANID"
      ASSERTED_RESULT_PARAM = "RESULT"

      SAMPLE_ID = "1"

      SAMPLE_RESULT_PARAMS = [
        {
          :result => {
            :actual => "Successful.",
            :expected => "successful"
          }
        },
        {
          :result => {
            :actual => "Failed - already charge during last 24H.",
            :expected => "errored"
          },
          :reason => {
            :expected => "already_charged"
          }
        },
        {
          :result => {
            :actual => "Failed - Do not has enough credit.",
            :expected => "failed"
          },
          :reason => {
            :expected => "not_enough_credit"
          }
        },
        {
          :result => {
            :actual => "Failed - MSISDN is not activated.",
            :expected => "errored"
          },
          :reason => {
            :expected => "number_not_activated"
          }
        },
        {
          :result => {
            :actual => "Failed - Invalid MSISDN.",
            :expected => "errored"
          },
          :reason => {
            :expected => "invalid_number"
          }
        },
        {
          :result => {
            :actual => "Failed.",
            :expected => "errored"
          }
        }
      ]

      def sample_result_params(sample_result)
        { ASSERTED_TRANSACTION_ID_PARAM => SAMPLE_ID, ASSERTED_RESULT_PARAM => sample_result[:result][:actual] }
      end

      SAMPLE_RESULT_PARAMS.each do |sample_result|
        context_string = "{'#{ASSERTED_TRANSACTION_ID_PARAM}' => '#{SAMPLE_ID}', '#{ASSERTED_RESULT_PARAM}' => '#{sample_result[:result][:actual]}'}"
        context "#params => #{context_string}" do
          subject { Qb.new(sample_result_params(sample_result)) }

          before do
            do_background_task(:queue_only => true) { subject.save! }
          end

          expected_reason_string = sample_result[:reason] ? "'#{sample_result[:reason][:expected]}'" : "nil"

          it "should queue a job for updating the charge request with args: '#{SAMPLE_ID}', '#{sample_result[:result][:expected]}', '#{ASSERTED_OPERATOR}', #{expected_reason_string}" do
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
