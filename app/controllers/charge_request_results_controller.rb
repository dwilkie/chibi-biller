class ChargeRequestResultsController < ApplicationController
  def create
    user, password = ActionController::HttpAuthentication::Basic::user_name_and_password(request)
    result = "charge_request_result/#{user}".classify.constantize.new(params)
    result.save!
    render(:nothing => true, :status => :created)
  end
end
