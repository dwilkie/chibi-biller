class ChargeRequestResult::Beeline < ::ChargeRequestResult::Base
  DIAMETER_SUCCESSFUL = "2001"
  DIAMETER_USER_UNKNOWN = [CHIBI_ERRORED, "diameter_cca_client_error_5030"]
  DIAMETER_RATING_FAILED = [CHIBI_ERRORED, "diameter_cca_client_error_5031"]

  DIAMETER_REASON_VALUES = {
    "4010" => CHIBI_INVALID_NUMBER,
    "4011" => CHIBI_CHARGE_NOT_APPLICABLE,
    "4012" => CHIBI_NOT_ENOUGH_CREDIT,
    "5030" => DIAMETER_USER_UNKNOWN,
    "5031" => DIAMETER_RATING_FAILED,
  }

  private

  def parse_result
    self.id = parse_diameter_data(params[:session_id])
    diameter_result = parse_diameter_data(params[:result_code])
    if diameter_result == DIAMETER_SUCCESSFUL
      self.result = CHIBI_SUCCESSFUL
    else
      chibi_full_result = DIAMETER_REASON_VALUES[diameter_result] || []
      self.result = chibi_full_result[0] || CHIBI_ERRORED
      self.reason = chibi_full_result[1]
    end
  end

  def parse_diameter_data(raw_data)
    (raw_data || "") =~ /\{(\d+)\}/
    ($~ || [])[1]
  end
end
