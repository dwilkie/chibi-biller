class ApplicationController < ActionController::Base
  # Prevent CSRF attacks by raising an exception.
  # For APIs, you may want to use :null_session instead.
  protect_from_forgery

  before_filter :authorize

  private

  def authorize
    authenticate_or_request_with_http_basic {|u, p| u == ENV["CHIBI_BILLER_CHARGE_REQUEST_USER_QB"] && p == ENV["CHIBI_BILLER_CHARGE_REQUEST_PASSWORD_QB"] }
  end
end
