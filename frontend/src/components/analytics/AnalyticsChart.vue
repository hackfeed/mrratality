<script>
import { Bar } from "vue3-chart-v2";

export default {
  mixins: [Bar],
  props: ["data", "options"],
  mounted() {
    this.renderChart(this.data, this.options);
  },
};
</script>
