<template>
  <section>
    <h2>Uploaded files</h2>
    <ul>
      <div @click="chooseFile(file)" class="entry" v-for="file in files" :key="file">
        <p class="filename">{{ file.name }}</p>
        <p class="uploaded-at">Uploaded at {{ parseDate(file.uploaded_at) }}</p>
      </div>
    </ul>
    <base-button @click="uploadNew(true)">Upload new</base-button>
  </section>
</template>

<script>
export default {
  props: ["files"],
  emits: ["upload-new", "choose-file", "is-uploaded"],
  methods: {
    parseDate(date) {
      const unixTime = Date.parse(date);
      const hsm = new Date(unixTime).toLocaleTimeString("ru-RU");
      const dt = new Date(unixTime).toLocaleDateString("ru-RU");
      return `${hsm} ${dt}`;
    },
    uploadNew(toUpload) {
      this.$emit("upload-new", toUpload);
    },
    chooseFile(file) {
      this.$emit("choose-file", file.name);
      this.$emit("is-uploaded", true);
    },
  },
};
</script>

<style scoped>
ul {
  padding: 0;
}
.entry {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 0 0 0.5rem;
  padding: 0 0.5rem;
  border: 1px rgba(0, 0, 0, 0.26) solid;
  border-radius: 0.5rem;
}
.entry:hover {
  cursor: pointer;
  border-color: none;
  background: rgb(231, 231, 231);
}
.filename {
  color: #389948;
  font-weight: bold;
}
.uploaded-at {
  color: rgb(185, 184, 184);
}
</style>
