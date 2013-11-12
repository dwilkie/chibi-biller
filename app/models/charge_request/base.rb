module ChargeRequest
  class Base
    require 'httparty'
    attr_accessor :transaction_id, :mobile_number

    def initialize(transaction_id, mobile_number)
      self.transaction_id = transaction_id
      self.mobile_number = mobile_number
    end

    def save!
      HTTParty.send(
        http_method,
        request_url,
        request_params
      )
    end
  end
end
