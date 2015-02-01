require 'rails_helper'

describe "Charge Request Result" do
  def authentication_params(options = {})
    {
      'HTTP_AUTHORIZATION' => ActionController::HttpAuthentication::Basic.encode_credentials(
        options[:user] || Rails.application.secrets[:charge_request_result_user_qb],
        options[:password] || Rails.application.secrets[:charge_request_result_password_qb]
      ),
      'REMOTE_ADDR' => options[:remote_ip] || Rails.application.secrets[:acl_whitelist].split(";").first
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
        expect(response.code).to eq("401")
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
          expect(response.code).to eq("201")
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
          expect(response).to be_nil
        end
      end
    end
  end
end
