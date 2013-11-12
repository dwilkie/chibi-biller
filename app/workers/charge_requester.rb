class ChargeRequester
  @queue = :charge_requester_queue

  def self.perform(charge_request_id, operator, mobile_number)
    request = ::ChargeRequest::Base.initialize_by_operator(operator, charge_request_id, mobile_number)
    request.save!
  end
end
