# this job will be queued by another app, therefore only perform method is important
class BeelineChargeRequestUpdaterJob < ActiveJob::Base
  def perform(session_id, result_code)
    ::ChargeRequestResult::Beeline.new(:session_id => session_id, :result_code => result_code).save!
  end
end
