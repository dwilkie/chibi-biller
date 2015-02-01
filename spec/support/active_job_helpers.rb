module ActiveJobHelpers
  include ActiveJob::TestHelper

  def trigger_job(options = {}, &block)
    if options.delete(:queue_only)
      yield
    else
      perform_enqueued_jobs { yield }
    end
  end

  def chibi_charge_request_updater_job
    enqueued_jobs.first
  end

  def assert_chibi_charge_request_updater_job(id, result, operator, reason = nil)
    expect(chibi_charge_request_updater_job).not_to be_nil
    expect(chibi_charge_request_updater_job[:args]).to eq([id, result, operator, reason])
  end
end
