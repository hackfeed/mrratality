<template>
  <div>
    <base-dialog :show="!!error" title="An error occured!" @close="handleError">
      <p>{{ error }}</p>
    </base-dialog>
    <base-card class="centered">
      <div v-if="!dataIsUploaded && !isLoading">
        <h2>Please upload CSV file below</h2>
        <analytics-form @upload-data="uploadData"></analytics-form>
      </div>
      <div v-else-if="isLoading">
        <h2>Data uploading is in process</h2>
        <base-spinner></base-spinner>
      </div>
      <div v-else-if="dataIsUploaded">
        <h2>Data is uploaded. Choose periods for MRR report</h2>
        <analytics-periods-form @load-data="loadData"></analytics-periods-form>
      </div>
    </base-card>
    <base-card v-if="dataIsLoaded" class="report-card">
      <analytics-chart :data="analyticsData" :options="analyticsOptions"></analytics-chart>
    </base-card>
  </div>
</template>

<script>
import AnalyticsForm from "../components/analytics/AnalyticsForm.vue";
import AnalyticsPeriodsForm from "../components/analytics/AnalyticsPeriodsForm.vue";
import AnalyticsChart from "../components/analytics/AnalyticsChart.vue";
export default {
  components: { AnalyticsForm, AnalyticsPeriodsForm, AnalyticsChart },
  data() {
    return {
      isLoading: false,
      loadError: false,
      error: null,
    };
  },
  computed: {
    dataIsUploaded() {
      return this.$store.getters["analytics/dataIsUploaded"];
    },
    dataIsLoaded() {
      return this.$store.getters["analytics/dataIsLoaded"];
    },
    analyticsData() {
      return this.$store.getters["analytics/data"];
    },
    analyticsOptions() {
      return this.$store.getters["analytics/dataOptions"];
    },
  },
  methods: {
    async uploadData(data) {
      this.isLoading = true;
      try {
        await this.$store.dispatch("analytics/uploadData", data);
      } catch (error) {
        this.error = error.message || "Something went wrong!";
      }
      this.isLoading = false;
    },
    async loadData(data) {
      this.isLoading = true;
      try {
        await this.$store.dispatch("analytics/loadData", data);
      } catch (error) {
        this.error = error.message || "Something went wrong!";
      }
      this.isLoading = false;
    },
    handleError() {
      this.error = null;
    },
  },
};
</script>

<style scoped>
.centered {
  text-align: center;
}
</style>
