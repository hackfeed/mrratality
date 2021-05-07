import { createRouter, createWebHistory } from "vue-router";

import Analytics from "./pages/Analytics.vue";
import Index from "./pages/Index.vue";
import NotFound from "./pages/NotFound.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", component: Index },
    { path: "/analytics", component: Analytics },
    { path: "/:notFound(.*)", component: NotFound },
  ],
});

export default router;
