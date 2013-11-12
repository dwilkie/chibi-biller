module OperatorMethods
  extend ActiveSupport::Concern

  module ClassMethods
    def operator
      self.name.demodulize.underscore
    end

    def initialize_by_operator(operator, *args)
      "#{self.name.deconstantize.underscore}/#{operator.to_s.underscore}".classify.constantize.new(*args)
    end
  end
end
