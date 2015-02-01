module ChargeRequest
  class Beeline < ::ChargeRequest::Base
    def save!
      Resque::Job.create(
        Rails.application.secrets[:beeline_charge_request_queue],
        Rails.application.secrets[:beeline_charge_request_worker],
        transaction_id.to_s, mobile_number
      )
    end
  end
end
