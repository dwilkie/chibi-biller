module ChargeRequestResult
  class Base
    attr_accessor :params

    def save!
      Resque::Job.create(
        ENV["CHIBI_CHARGE_REQUEST_UPDATER_QUEUE"],
        ENV["CHIBI_CHARGE_REQUEST_UPDATER_WORKER"],
        id, result, operator, reason
      )
    end
  end
end
