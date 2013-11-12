module ChargeRequestResult
  class Base
    include OperatorMethods

    attr_accessor :params, :id, :result, :reason

    def initialize(params = {})
      self.params = params
    end

    def save!
      Resque::Job.create(
        ENV["CHIBI_CHARGE_REQUEST_UPDATER_QUEUE"],
        ENV["CHIBI_CHARGE_REQUEST_UPDATER_WORKER"],
        id, result, self.class.operator, reason
      )
    end

    def error(reason = nil)
      self.result = "errored"
      self.reason = reason
    end
  end
end
