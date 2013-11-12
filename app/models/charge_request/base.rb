module ChargeRequest
  class Base
    include OperatorMethods

    attr_accessor :transaction_id, :mobile_number

    def initialize(transaction_id, mobile_number)
      self.transaction_id = transaction_id
      self.mobile_number = mobile_number
    end

    def save!
      client_result = HTTParty.send(http_method, request_url, request_params)
      unless client_result.success?
        charge_request_result = ::ChargeRequestResult::Base.initialize_by_operator(self.class.operator)
        charge_request_result.id = transaction_id
        charge_request_result.error("Client Error: #{client_result.message}")
        charge_request_result.save!
      end
    end
  end
end
