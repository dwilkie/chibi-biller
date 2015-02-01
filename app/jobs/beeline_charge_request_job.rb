# this job will be run by another app therefore only the queue name is important
class BeelineChargeRequestJob < ActiveJob::Base
  queue_as :beeline_charge_request_queue
end
