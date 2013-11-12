module ChargeRequestResultExamples
  shared_examples_for "a charge request result" do
    describe "#error(reason)" do
      it "should update #result and #reason" do
        subject.error("reason")
        subject.result.should == "errored"
        subject.reason.should == "reason"
      end
    end

    [:id, :result, :reason, :params].each do |accessor|
      describe "##{accessor}" do
        it "should be an accessor" do
          subject.send("#{accessor}=", accessor)
          subject.send(accessor).should == accessor
        end
      end
    end
  end
end
