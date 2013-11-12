class ChargeRequestResultsController < ApplicationController
  def create
    user, password = ActionController::HttpAuthentication::Basic::user_name_and_password(request)
    request = "charge_request_result/#{user}".classify.constantize.new(params)
    request.save!
    render(:nothing => true, :status => :created)
  end
end
