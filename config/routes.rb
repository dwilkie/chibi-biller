ChibiBiller::Application.routes.draw do
  request_ip_constraint = lambda { |request| request.remote_ip == ENV["QB_BILLING_API_HOST_IP"] }
  constraints request_ip_constraint do
    resources :charge_request_results, :only => :create
  end
end
