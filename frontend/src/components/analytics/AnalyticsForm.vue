<template>
  <section>
    <div>
      <label for="file" class="flat">Browse</label>
      <input type="file" id="file" ref="file" @change="handleFileUpload" />
    </div>
    <base-button @click="submitFile">Upload</base-button>
  </section>
</template>

<script>
export default {
  emits: ["upload-data"],
  data() {
    return {
      file: "",
    };
  },
  methods: {
    handleFileUpload() {
      this.file = this.$refs.file.files[0];
    },
    submitFile() {
      let formData = new FormData();
      formData.append("file", this.file);
      this.$emit("upload-data", formData);
    },
  },
};
</script>

<style scoped>
label {
  margin-bottom: 1rem;
}
#file {
  opacity: 0;
  position: absolute;
  z-index: -1;
}
.flat {
  text-decoration: none;
  padding: 0.75rem 1.5rem;
  font: inherit;
  background-color: transparent;
  color: #389948;
  border: none;
  cursor: pointer;
  border-radius: 30px;
  margin-right: 0.5rem;
  display: inline-block;
}
.flat:hover,
.flat:active {
  background-color: #70b87c;
  color: white;
}
</style>
