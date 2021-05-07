import { createApp, defineAsyncComponent } from "vue";

import App from "./App.vue";
import BaseButton from "./components/ui/BaseButton.vue";
import BaseCard from "./components/ui/BaseCard.vue";
import BaseSpinner from "./components/ui/BaseSpinner.vue";
import router from "./router.js";
import store from "./store/index.js";

const BaseDialog = defineAsyncComponent(() => {
  import("./components/ui/BaseDialog.vue");
});

const app = createApp(App);

app.use(router);
app.use(store);

app.component("base-card", BaseCard);
app.component("base-button", BaseButton);
app.component("base-spinner", BaseSpinner);
app.component("base-dialog", BaseDialog);

app.mount("#app");
