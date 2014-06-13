class BeelineChargeRequestUpdater
  @queue = :beeline_charge_request_updater_queue

  def self.perform(session_id, result_code)
    result = ::ChargeRequestResult::Beeline.new(:session_id => session_id, :result_code => result_code)
    result.save!
  end
end
