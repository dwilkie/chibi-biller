module ChargeRequestResultExamples
  shared_examples_for "a charge request result" do
    describe "#error(reason)" do
      it "should update #result and #reason" do
        subject.error("reason")
        expect(subject.result).to eq("errored")
        expect(subject.reason).to eq("reason")
      end
    end

    [:id, :result, :reason, :params].each do |accessor|
      describe "##{accessor}" do
        it "should be an accessor" do
          subject.send("#{accessor}=", accessor)
          expect(subject.send(accessor)).to eq(accessor)
        end
      end
    end
  end
end
