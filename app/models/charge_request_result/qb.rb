module ChargeRequestResult
  class Qb < Base
    QB_SUCCESSFUL = "successful"
    QB_TRANSACTION_ID_PARAM = "TRANID"
    QB_RESULT_PARAM = "RESULT"
    QB_RESULT_REASON_DELIMITER = "-"

    QB_REASON_VALUES = {
      "do not has enough credit"        => CHIBI_NOT_ENOUGH_CREDIT,
      "already charge during last 24h"  => CHIBI_ALREADY_CHARGED,
      "msisdn is not activated"         => CHIBI_NUMBER_NOT_ACTIVATED,
      "invalid msisdn"                  => CHIBI_INVALID_NUMBER
    }

    private

    def parse_result
      self.id = params[QB_TRANSACTION_ID_PARAM]
      qb_full_result = sanitized(params[QB_RESULT_PARAM]).split(QB_RESULT_REASON_DELIMITER)
      qb_result = sanitized(qb_full_result[0])
      qb_reason = sanitized(qb_full_result[1])
      if qb_result == QB_SUCCESSFUL
        self.result = CHIBI_SUCCESSFUL
      else
        chibi_full_result = QB_REASON_VALUES[qb_reason] || []
        self.result = chibi_full_result[0] || CHIBI_ERRORED
        self.reason = chibi_full_result[1]
      end
    end

    def sanitized(string)
      string.to_s.strip.downcase.gsub(/\.$/, "")
    end
  end
end
