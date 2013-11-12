module ChargeRequestExamples
  shared_examples_for "a charge request" do
    [:transaction_id, :mobile_number].each do |accessor|
      describe "##{accessor}" do
        it "should be an accessor" do
          subject.send("#{accessor}=", accessor)
          subject.send(accessor).should == accessor
        end
      end
    end
  end
end
