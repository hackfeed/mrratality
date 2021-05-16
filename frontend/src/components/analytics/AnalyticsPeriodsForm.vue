<template>
  <section>
    <div class="periods">
      <div class="period">
        <label for="periodStart">Period Start</label>
        <input type="month" name="periodStart" id="periodStart" v-model="periodStart" />
      </div>
      <div class="period">
        <label for="periodEnd">Period End</label>
        <input type="month" name="periodEnd" id="periodEnd" v-model="periodEnd" />
      </div>
    </div>
    <div>
      <base-button mode="flat" @click="uploadNew(false)">Back to files</base-button>
      <base-button @click="submitPeriods">Load report</base-button>
    </div>
  </section>
</template>

<script>
export default {
  emits: ["load-data", "upload-new"],
  data() {
    return {
      periodStart: "2021-01",
      periodEnd: "2021-01",
    };
  },
  methods: {
    submitPeriods() {
      this.$emit("load-data", { periodStart: this.periodStart, periodEnd: this.periodEnd });
    },
    uploadNew(toUpload) {
      this.$emit("upload-new", toUpload);
    },
  },
};
</script>

<style scoped>
label {
  font-weight: bold;
  display: block;
  margin-bottom: 0.5rem;
}
input[type="month"] {
  border: 1px solid #389948;
  border-radius: 5px;
  color: black;
  text-align: center;
  font-family: "Roboto", sans-serif;
  cursor: pointer;
}
input[type="month"]:focus {
  outline: none;
}
.periods {
  display: flex;
}
.period {
  margin: 1rem auto;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
</style>
