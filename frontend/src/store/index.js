import analyticsModule from "./modules/analytics/index.js";
import authModule from "./modules/auth/index.js";
import { createStore } from "vuex";

const store = createStore({
  modules: {
    auth: authModule,
    analytics: analyticsModule,
  },
});

export default store;
