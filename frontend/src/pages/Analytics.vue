<template>
  <div>
    <base-dialog :show="!!error" title="An error occured!" @close="handleError">
      <p>{{ error }}</p>
    </base-dialog>
    <base-card class="centered">
      <div v-if="isLoading && !uploadNew && !isUploaded && !isLoaded">
        <h2>Loading user files</h2>
        <base-spinner></base-spinner>
      </div>
      <div v-else-if="!isLoading && !uploadNew && !isUploaded">
        <analytics-files
          :files="analyticsFiles"
          @upload-new="setUploadNew"
          @choose-file="setFile"
          @is-uploaded="setIsUploaded"
        ></analytics-files>
      </div>
      <div v-if="uploadNew && !isUploaded && !isLoading">
        <h2>Please upload CSV file below</h2>
        <analytics-form @upload-data="uploadData" @upload-new="setUploadNew"></analytics-form>
      </div>
      <div v-else-if="isLoading">
        <h2>Data uploading is in process</h2>
        <base-spinner></base-spinner>
      </div>
      <div v-else-if="isUploaded">
        <h2>Choose periods for MRR report</h2>
        <analytics-periods-form
          @load-data="loadData"
          @upload-new="setUploadNew"
        ></analytics-periods-form>
      </div>
    </base-card>
    <base-card v-if="isLoaded" class="report-card">
      <analytics-chart :data="analyticsData" :options="analyticsOptions"></analytics-chart>
    </base-card>
  </div>
</template>

<script>
import AnalyticsFiles from "../components/analytics/AnalyticsFiles.vue";
import AnalyticsForm from "../components/analytics/AnalyticsForm.vue";
import AnalyticsPeriodsForm from "../components/analytics/AnalyticsPeriodsForm.vue";
import AnalyticsChart from "../components/analytics/AnalyticsChart.vue";
export default {
  components: { AnalyticsFiles, AnalyticsForm, AnalyticsPeriodsForm, AnalyticsChart },
  data() {
    return {
      uploadNew: false,
      isUploaded: false,
      isLoading: false,
      isLoaded: false,
      loadError: false,
      file: null,
      error: null,
    };
  },
  computed: {
    analyticsFiles() {
      return this.$store.getters["analytics/files"];
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
      this.isUploaded = true;
      this.file = this.$store.getters["analytics/file"];
    },
    async loadData(data) {
      this.isLoading = true;
      try {
        await this.$store.dispatch("analytics/loadData", data);
      } catch (error) {
        this.error = error.message || "Something went wrong!";
      }
      this.isLoading = false;
      this.isLoaded = true;
    },
    async loadFiles() {
      this.isLoading = true;
      try {
        await this.$store.dispatch("analytics/loadFiles");
      } catch (error) {
        this.error = error.message || "Something went wrong!";
      }
      this.isLoading = false;
    },
    setUploadNew(data) {
      if (data === false) {
        this.isUploaded = false;
        this.isLoaded = false;
      }
      this.uploadNew = data;
    },
    setFile(data) {
      this.file = data;
    },
    setIsUploaded(data) {
      this.isUploaded = data;
    },
    handleError() {
      this.error = null;
    },
  },
  created() {
    this.loadFiles();
  },
};
</script>

<style scoped>
.centered {
  text-align: center;
}
</style>
