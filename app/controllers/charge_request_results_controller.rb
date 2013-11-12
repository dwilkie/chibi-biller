class ChargeRequestResultsController < ApplicationController
  def create
    user, password = ActionController::HttpAuthentication::Basic::user_name_and_password(request)
    result = ChargeRequestResult::Base.initialize_by_operator(user, params)
    result.save!
    render(:nothing => true, :status => :created)
  end
end
