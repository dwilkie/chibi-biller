class BeelineChargeRequestJob < ActiveJob::Base
  queue_as Rails.application.secrets[:beeline_charge_request_queue]
end
