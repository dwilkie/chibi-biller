module ChargeRequestHelpers
  private

  def asserted_url(operator, options = {})
    url = URI.parse(Rails.application.secrets[:"charge_request_url_#{operator}"])
    url.query = Rack::Utils.build_query(options[:query]) if options[:query]
    url.to_s
  end

  def expect_charge_request(options = {}, &block)
    response_type = options.delete(:response_type) || :success
    cassette = options.delete(:cassette) || [options.delete(:operator), response_type].join("_")
    VCR.use_cassette(cassette, :erb => options) { yield }
  end
end
