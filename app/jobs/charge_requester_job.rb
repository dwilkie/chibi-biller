class ChargeRequesterJob < ActiveJob::Base
  def perform(charge_request_id, operator, mobile_number)
    ::ChargeRequest::Base.initialize_by_operator(operator, charge_request_id, mobile_number).save!
  end
end
