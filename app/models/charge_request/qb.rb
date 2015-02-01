module ChargeRequest
  class Qb < ::ChargeRequest::Base
    PARAM_KEYS = {
      :mobile_number => "MSISDN",
      :transaction_id => "TRANID"
    }

    def save!
      client_result = HTTParty.send(http_method, request_url, request_params)
      unless client_result.success?
        charge_request_result = ::ChargeRequestResult::Base.initialize_by_operator(self.class.operator)
        charge_request_result.id = transaction_id
        charge_request_result.error("Client Error: #{client_result.message}")
        charge_request_result.save!
      end
    end

    private

    def http_method
      :get
    end

    def request_url
      Rails.application.secrets[:charge_request_url_qb]
    end

    def request_params
      {:query => {PARAM_KEYS[:transaction_id] => transaction_id, PARAM_KEYS[:mobile_number] => mobile_number}}
    end
  end
end
