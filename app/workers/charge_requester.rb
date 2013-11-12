class ChargeRequester
  @queue = :charge_requester_queue

  def self.perform(charge_request_id, operator, mobile_number)
    request = "charge_request/#{operator}".classify.constantize.new(charge_request_id, mobile_number)
    request.save!
  end
end
