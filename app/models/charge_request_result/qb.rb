module ChargeRequestResult
  class Qb < Base
    RESULT_VALUES = {
      :successful => "Successful.",
      :failed => "Failed."
    }

    PARAM_KEYS = {
      :id => "TRANID",
      :result => "RESULT",
      :reason => "REASON"
    }

    def initialize(params)
      self.params = params
    end

    private

    def id
      params[PARAM_KEYS[:id]]
    end

    def result
      if successful?
        "successful"
      elsif failed?
        "failed"
      else
        "errored"
      end
    end

    def operator
      self.class.name.demodulize.underscore
    end

    def reason
      params[PARAM_KEYS[:reason]]
    end

    def successful?
      params[PARAM_KEYS[:result]] == RESULT_VALUES[:successful]
    end

    def failed?
      params[PARAM_KEYS[:result]] == RESULT_VALUES[:failed]
    end
  end
end
