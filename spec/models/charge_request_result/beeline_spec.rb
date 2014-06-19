require 'spec_helper'

describe ChargeRequestResult::Beeline do
  include OperatorExamples
  include ChargeRequestResultExamples

  BEELINE_ASSERTED_OPERATOR = "beeline"

  let(:asserted_operator) { BEELINE_ASSERTED_OPERATOR }
  let(:initialization_args) { [] }

  it_should_behave_like "an operator"
  it_should_behave_like "a charge request result"

  describe "#save!" do
    include ResqueHelpers

    BEELINE_ASSERTED_TRANSACTION_ID_PARAM = :session_id
    BEELINE_ASSERTED_RESULT_PARAM = :result_code

    BEELINE_SAMPLE_RESULT_PARAMS = [
      {
        :transaction_id => {
          :actual => "UTF8String{3319142568},Padding:2",
          :expected => "3319142568"
        },
        :result => {
          :actual => "Unsigned32{2001}",
          :expected => "successful"
        }
      },
      {
        :transaction_id => {
          :actual => "UTF8String{2},Padding:313",
          :expected => "2"
        },
        :result => {
          :actual => "Unsigned32{4010}",
          :expected => "errored"
        },
        :reason => {
          :expected => "invalid_number"
        }
      },
      {
        :transaction_id => {
          :actual => "UTF8String{3243},Padding:112",
          :expected => "3243"
        },
        :result => {
          :actual => "Unsigned32{4012}",
          :expected => "failed"
        },
        :reason => {
          :expected => "not_enough_credit"
        }
      },
      {
        :transaction_id => {
          :actual => "UTF8String{1234},Padding:112",
          :expected => "1234"
        },
        :result => {
          :actual => "Unsigned32{5031}",
          :expected => "errored"
        },
        :reason => {
          :expected => "diameter_cca_client_error_5031"
        }
      }
    ]

    def sample_result_params(sample_result)
      { BEELINE_ASSERTED_TRANSACTION_ID_PARAM => sample_result[:transaction_id][:actual], BEELINE_ASSERTED_RESULT_PARAM => sample_result[:result][:actual] }
    end

    BEELINE_SAMPLE_RESULT_PARAMS.each do |sample_result|
      context_string = "{'#{BEELINE_ASSERTED_TRANSACTION_ID_PARAM}' => '#{sample_result[:transaction_id][:actual]}', '#{BEELINE_ASSERTED_RESULT_PARAM}' => '#{sample_result[:result][:actual]}'}"
      context "#params => #{context_string}" do
        subject { described_class.new(sample_result_params(sample_result)) }

        before do
          do_background_task(:queue_only => true) { subject.save! }
        end

        expected_reason_string = sample_result[:reason] ? "'#{sample_result[:reason][:expected]}'" : "nil"

        it "should queue a job for updating the charge request with args: '#{sample_result[:transaction_id][:expected]}', '#{sample_result[:result][:expected]}', '#{BEELINE_ASSERTED_OPERATOR}', #{expected_reason_string}" do
          assert_chibi_charge_request_updater_job(
            sample_result[:transaction_id][:expected],
            sample_result[:result][:expected],
            asserted_operator,
            (sample_result[:reason] || {})[:expected]
          )
        end
      end
    end
  end
end
