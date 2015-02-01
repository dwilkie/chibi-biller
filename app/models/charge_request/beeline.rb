class ChargeRequest::Beeline  < ::ChargeRequest::Base
  def save!
    BeelineChargeRequestJob.perform_later(transaction_id.to_s, mobile_number)
  end
end
