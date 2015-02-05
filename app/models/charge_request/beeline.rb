class ChargeRequest::Beeline  < ::ChargeRequest::Base
  def save!
    BeelineChargeRequestJob.perform_later(
      transaction_id.to_s,
      mobile_number,
      BeelineChargeRequestUpdaterJob.queue_name,
      BeelineChargeRequestUpdaterJob.to_s,
      Rails.application.secrets[:beeline_billing_server_address]
    )
  end
end
