require 'spec_helper'

describe "Charge Request Result" do
  def authentication_params(options = {})
    {
      'HTTP_AUTHORIZATION' => ActionController::HttpAuthentication::Basic.encode_credentials(
        options[:user] || ENV["CHIBI_BILLER_CHARGE_REQUEST_USER_QB"],
        options[:password] || ENV["CHIBI_BILLER_CHARGE_REQUEST_PASSWORD_QB"]
      ),
      'REMOTE_ADDR' => options[:remote_ip] || ENV["QB_BILLING_API_HOST_IP"]
    }
  end

  def post_charge_request_results_path(options = {})
    post(
      charge_request_results_path, {}, authentication_params(options)
    )
  end

  describe "POST /charge_request_results" do
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
        before do
          post_charge_request_results_path
        end

        it "should grant access" do
          response.code.should == "201"
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
