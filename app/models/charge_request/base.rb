class ChargeRequest::Base
  include OperatorMethods

  attr_accessor :transaction_id, :mobile_number

  def initialize(transaction_id, mobile_number)
    self.transaction_id = transaction_id
    self.mobile_number = mobile_number
  end
end
