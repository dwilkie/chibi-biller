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

  def assert_chibi_charge_request_updater_job(id, result, operator, reason = nil)
    job = ResqueSpec.queues[ENV["CHIBI_CHARGE_REQUEST_UPDATER_QUEUE"]].first
    job.should_not be_nil
    job[:class].should == ENV["CHIBI_CHARGE_REQUEST_UPDATER_WORKER"]
    job[:args].should == [id, result, operator, reason]
  end
end
