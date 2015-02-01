module ResqueHelpers
  EXTERNAL_QUEUES = [Rails.application.secrets[:chibi_charge_request_updater_queue]]

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
    @chibi_charge_request_updater_job ||= ResqueSpec.queues[Rails.application.secrets[:chibi_charge_request_updater_queue]].first
  end

  def assert_chibi_charge_request_updater_job(id, result, operator, reason = nil)
    expect(chibi_charge_request_updater_job).not_to be_nil
    expect(chibi_charge_request_updater_job[:class]).to eq(Rails.application.secrets[:chibi_charge_request_updater_worker])
    expect(chibi_charge_request_updater_job[:args]).to eq([id, result, operator, reason])
  end
end
