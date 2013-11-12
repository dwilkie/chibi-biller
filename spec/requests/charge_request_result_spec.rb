require 'spec_helper'

describe "Charge Request Result" do
  def authentication_params(options = {})
    {
      'HTTP_AUTHORIZATION' => ActionController::HttpAuthentication::Basic.encode_credentials(
        options[:user] || ENV["CHARGE_REQUEST_RESULT_USER_QB"],
        options[:password] || ENV["CHARGE_REQUEST_RESULT_PASSWORD_QB"]
      ),
      'REMOTE_ADDR' => options[:remote_ip] || ENV["ACL_WHITELIST"].split(";").first
    }
  end

  def post_charge_request_results_path(options = {})
    params = options.slice!(:user, :password, :remote_ip)
    post(
      charge_request_results_path, params, authentication_params(options)
    )
  end

  describe "POST '/charge_request_results'" do
    context "with incorrect authentication" do
      before do
        post_charge_request_results_path(:user => "wrong", :password => "bad")
      end

      it "should prevent access" do
        response.code.should == "401"
      end
    end

    context "with correct authentication" do
      context "and an allowed ip" do
        include ResqueHelpers

        before do
          do_background_task(:queue_only => true) do
            post_charge_request_results_path(
              "TRANID" => "1",
              "RESULT" => "Successful."
            )
          end
        end

        it "should queue a job for updating the charge request" do
          response.code.should == "201"
          assert_chibi_charge_request_updater_job("1", "successful", "qb", nil)
        end
      end

      context "and a disallowed ip" do
        before do
          expect {
            post_charge_request_results_path(:remote_ip => "123.456.789.1")
          }.to raise_exception(ActionController::RoutingError)
        end

        it "should deny me access" do
          response.should be_nil
        end
      end
    end
  end
end
