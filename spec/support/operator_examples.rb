module OperatorExamples
  shared_examples_for "an operator" do
    describe ".operator" do
      it "should return the operator name" do
        expect(subject.class.operator).to eq(asserted_operator)
      end
    end

    describe ".initialize_by_operator(operator, *args)" do
      it "should return a new instance of the operator's request or result" do
        result = subject.class.initialize_by_operator(asserted_operator, *initialization_args)
        expect(result).to be_a(subject.class)
      end
    end
  end
end
