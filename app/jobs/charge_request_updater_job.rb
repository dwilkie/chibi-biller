# this job will be run on Chibi therefore only the queue name is important
class ChargeRequestUpdaterJob < ActiveJob::Base
  queue_as :critical
end
