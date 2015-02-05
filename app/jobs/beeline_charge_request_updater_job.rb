# this job will be queued by another app, therefore only perform method is important
class BeelineChargeRequestUpdaterJob < ActiveJob::Base
  queue_as :beeline_charge_request_updater_queue

  attr_accessor :jid

  def perform(session_id, result_code)
    ::ChargeRequestResult::Beeline.new(:session_id => session_id, :result_code => result_code).save!
  end
end
