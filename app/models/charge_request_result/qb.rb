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

    def id
      super || params[PARAM_KEYS[:id]]
    end

    def result
      result = super
      return result if result
      if successful?
        "successful"
      elsif failed?
        "failed"
      else
        "errored"
      end
    end

    def reason
      super || params[PARAM_KEYS[:reason]]
    end

    private

    def successful?
      params[PARAM_KEYS[:result]] == RESULT_VALUES[:successful]
    end

    def failed?
      params[PARAM_KEYS[:result]] == RESULT_VALUES[:failed]
    end
  end
end
