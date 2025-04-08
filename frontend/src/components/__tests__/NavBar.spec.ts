import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createRouter, createWebHistory } from "vue-router";
import { createPinia, setActivePinia } from "pinia";
import NavBar from "../NavBar.vue";
import { nextTick } from "vue";

// Mock the store
vi.mock("../stores/stockStore", () => ({
  useStockStore: () => ({
    syncStocks: vi.fn().mockResolvedValue(undefined),
  }),
}));

// Create mock router
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", component: { template: "<div>Home</div>" } },
    { path: "/recommendations", component: { template: "<div>Recommendations</div>" } },
  ],
});

describe("NavBar", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("renders properly", () => {
    const wrapper = mount(NavBar, {
      global: {
        plugins: [router],
      },
    });

    expect(wrapper.text()).toContain("StockAnalyzer");
    expect(wrapper.text()).toContain("Stocks");
    expect(wrapper.text()).toContain("Recommendations");
    expect(wrapper.text()).toContain("Sync Stocks");
  });

  it("toggles mobile menu when burger icon is clicked", async () => {
    const wrapper = mount(NavBar, {
      global: {
        plugins: [router],
      },
    });

    // Initially the mobile menu should be hidden
    expect(wrapper.find(".md\\:hidden > div").exists()).toBe(false);

    // Click the burger icon
    await wrapper.find(".md\\:hidden button").trigger("click");

    // The mobile menu should now be visible
    expect(wrapper.find(".md\\:hidden > div").exists()).toBe(true);

    // Click the X icon to close
    await wrapper.find(".md\\:hidden button").trigger("click");

    // The mobile menu should be hidden again
    expect(wrapper.find(".md\\:hidden > div").exists()).toBe(false);
  });

  it("closes mobile menu when navigation links are clicked", async () => {
    const wrapper = mount(NavBar, {
      global: {
        plugins: [router],
      },
    });

    // Open the mobile menu
    await wrapper.find(".md\\:hidden button").trigger("click");
    expect(wrapper.find(".md\\:hidden > div").exists()).toBe(true);

    // Click the Stocks link
    const stocksLink = wrapper.findAll(".md\\:hidden a")[0];
    await stocksLink.trigger("click");

    // The mobile menu should be closed
    expect(wrapper.find(".md\\:hidden > div").exists()).toBe(false);
  });
});
