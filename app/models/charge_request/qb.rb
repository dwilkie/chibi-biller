module ChargeRequest
  class Qb < Base
    PARAM_KEYS = {
      :mobile_number => "MSISDN",
      :transaction_id => "TRANID"
    }

    private

    def http_method
      :get
    end

    def request_url
      ENV["CHARGE_REQUEST_URL_QB"]
    end

    def request_params
      {:query => {PARAM_KEYS[:transaction_id] => transaction_id, PARAM_KEYS[:mobile_number] => mobile_number}}
    end
  end
end
