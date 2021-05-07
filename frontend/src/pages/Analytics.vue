<template>
  <base-card class="centered">
    <div v-if="!dataLoaded && !isLoading">
      <h2>Please upload CSV file below</h2>
      <analytics-form @save-data="saveData"></analytics-form>
    </div>
    <div v-else-if="isLoading">
      <h2>Data loading is in process</h2>
      <base-spinner></base-spinner>
    </div>
    <div v-else-if="dataLoaded">
      <h2>Data is loaded. Have fun!</h2>
    </div>
  </base-card>
</template>

<script>
import AnalyticsForm from "../components/analytics/AnalyticsForm.vue";
export default {
  components: { AnalyticsForm },
  data() {
    return {
      dataLoaded: false,
      isLoading: false,
    };
  },
  methods: {
    saveData(data) {
      //   this.$store.dispatch("analytics/uploadData", data);
      console.log(data);
      this.isLoading = true;
      setTimeout(() => {
        this.isLoading = false;
        this.dataLoaded = true;
      }, 2000);
    },
  },
};
</script>

<style scoped>
.centered {
  text-align: center;
}
</style>
