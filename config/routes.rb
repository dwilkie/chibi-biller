ChibiBiller::Application.routes.draw do
  request_ip_constraint = lambda { |request| Rails.application.secrets[:acl_whitelist].to_s.split(";").include?(request.remote_ip) }
  constraints request_ip_constraint do
    resources :charge_request_results, :only => :create
  end
end
