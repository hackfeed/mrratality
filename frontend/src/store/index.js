import analyticsModule from "./modules/analytics/index.js";
import { createStore } from "vuex";

const store = createStore({
  modules: {
    analytics: analyticsModule,
  },
});

export default store;
