module ChargeRequest
  class Beeline < ::ChargeRequest::Base
    def save!
      Resque::Job.create(
        ENV["BEELINE_CHARGE_REQUEST_QUEUE"],
        ENV["BEELINE_CHARGE_REQUEST_WORKER"],
        transaction_id.to_s, mobile_number
      )
    end
  end
end