module ChargeRequestResult
  class Base
    include OperatorMethods

    CHIBI_SUCCESSFUL = "successful"
    CHIBI_FAILED = "failed"
    CHIBI_ERRORED = "errored"

    CHIBI_NOT_ENOUGH_CREDIT    = [CHIBI_FAILED,  "not_enough_credit"]
    CHIBI_ALREADY_CHARGED      = [CHIBI_ERRORED, "already_charged"]
    CHIBI_NUMBER_NOT_ACTIVATED = [CHIBI_ERRORED, "number_not_activated"]
    CHIBI_INVALID_NUMBER       = [CHIBI_ERRORED, "invalid_number"]

    attr_accessor :params, :id, :result, :reason

    def initialize(params = {})
      self.params = params
      parse_result
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
