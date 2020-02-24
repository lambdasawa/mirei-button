import { shallowMount } from "@vue/test-utils";
import Trimmer from "@/components/Trimmer.vue";

describe("Trimmer.vue", () => {
  it("renders props.msg when passed", () => {
    const msg = "new message";
    const wrapper = shallowMount(Trimmer, {
      propsData: { msg }
    });
    expect(wrapper.text()).toMatch(msg);
  });
});
