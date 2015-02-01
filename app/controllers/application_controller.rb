class ApplicationController < ActionController::Base
  # Prevent CSRF attacks by raising an exception.
  # For APIs, you may want to use :null_session instead.
  protect_from_forgery

  before_filter :authorize

  private

  def authorize
    authenticate_or_request_with_http_basic do |u, p|
      u == Rails.application.secrets[:charge_request_result_user_qb] && p == Rails.application.secrets[:charge_request_result_password_qb]
    end
  end
end
