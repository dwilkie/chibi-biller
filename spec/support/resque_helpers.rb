module ResqueHelpers
  EXTERNAL_QUEUES = [ENV["CHIBI_CHARGE_REQUEST_UPDATER_QUEUE"]]

  private

  def do_background_task(options = {}, &block)
    ResqueSpec.reset!

    yield

    unless options.delete(:queue_only)
      # don't perform jobs that are not intended to be performed by this app
      queues = ResqueSpec.queues.reject { |k, v| EXTERNAL_QUEUES.include?(k.to_sym) }
      queues.keys.each do |queue_name|
        ResqueSpec.perform_all(queue_name)
      end
    end
  end

  def perform_background_job(queue_name)
    ResqueSpec.perform_all(queue_name)
  end

  def chibi_charge_request_updater_job
    @chibi_charge_request_updater_job ||= ResqueSpec.queues[ENV["CHIBI_CHARGE_REQUEST_UPDATER_QUEUE"]].first
  end

  def assert_chibi_charge_request_updater_job(id, result, operator, reason = nil)
    chibi_charge_request_updater_job.should_not be_nil
    chibi_charge_request_updater_job[:class].should == ENV["CHIBI_CHARGE_REQUEST_UPDATER_WORKER"]
    chibi_charge_request_updater_job[:args].should == [id, result, operator, reason]
  end
end
