ChibiBiller::Application.routes.draw do
  request_ip_constraint = lambda { |request| request.remote_ip == ENV["ACL_WHITELIST"].to_s.split(";").first }
  constraints request_ip_constraint do
    resources :charge_request_results, :only => :create
  end
end
