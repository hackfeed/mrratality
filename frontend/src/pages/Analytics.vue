<template>
  <div>
    <base-dialog :show="!!error" title="An error occured!" @close="handleError">
      <p>{{ error }}</p>
    </base-dialog>
    <base-card class="centered">
      <div v-if="!uploadNew && !isUploaded && !isLoading && analyticsFiles.length > 0">
        <analytics-files
          :files="analyticsFiles"
          @upload-new="setUploadNew"
          @choose-file="setFile"
          @delete-file="deleteFile"
          @is-uploaded="setIsUploaded"
        ></analytics-files>
      </div>
      <div v-else-if="!isUploaded && !isLoading">
        <h2>Please upload CSV file below</h2>
        <analytics-form
          :files-not-empty="analyticsFilesEmpty"
          @upload-data="uploadData"
          @upload-new="setUploadNew"
        ></analytics-form>
      </div>
      <div v-else-if="isLoading">
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
    <div v-if="isLoaded" class="centered">
      <base-card>
        <analytics-chart :data="analyticsData" :options="analyticsOptions"></analytics-chart>
      </base-card>
      <base-card>
        <analytics-grid :grid="analyticsGrid"></analytics-grid>
      </base-card>
    </div>
  </div>
</template>

<script>
import AnalyticsFiles from "../components/analytics/AnalyticsFiles.vue";
import AnalyticsForm from "../components/analytics/AnalyticsForm.vue";
import AnalyticsPeriodsForm from "../components/analytics/AnalyticsPeriodsForm.vue";
import AnalyticsChart from "../components/analytics/AnalyticsChart.vue";
import AnalyticsGrid from "../components/analytics/AnalyticsGrid.vue";

export default {
  components: {
    AnalyticsFiles,
    AnalyticsForm,
    AnalyticsPeriodsForm,
    AnalyticsChart,
    AnalyticsGrid,
  },
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
    analyticsGrid() {
      return this.$store.getters["analytics/grid"];
    },
    analyticsOptions() {
      return this.$store.getters["analytics/dataOptions"];
    },
    analyticsFilesEmpty() {
      return this.analyticsFiles.length === 0;
    },
  },
  methods: {
    async uploadData(data) {
      this.isLoading = true;
      try {
        await this.$store.dispatch("analytics/uploadData", data);
        this.isUploaded = true;
        this.file = this.$store.getters["analytics/file"];
        await this.loadFiles();
      } catch (error) {
        this.error = error.message || "Something went wrong!";
      }
      this.isLoading = false;
    },
    async loadData(data) {
      this.isLoaded = false;
      this.isLoading = true;
      data = {
        filename: this.file,
        period_start: data.periodStart + "-01",
        period_end: data.periodEnd + "-01",
      };
      try {
        await this.$store.dispatch("analytics/loadData", data);
        this.isLoaded = true;
      } catch (error) {
        this.error = error.message || "Something went wrong!";
      }
      this.isLoading = false;
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
    async deleteFile(data) {
      try {
        await this.$store.dispatch("analytics/deleteFile", data);
        await this.loadFiles();
      } catch (error) {
        this.error = error.message || "Something went wrong!";
      }
    },
    setUploadNew(data) {
      if (data === false) {
        this.isUploaded = false;
        this.isLoaded = false;
        this.file = null;
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
