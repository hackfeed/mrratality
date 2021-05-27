<template>
  <h2>{{ grid.title }}</h2>
  <div id="wrapper"></div>
</template>

<script>
import { Grid } from "gridjs";
import "gridjs/dist/theme/mermaid.css";

export default {
  props: ["grid"],
  mounted() {
    new Grid({ columns: this.grid.cols, data: this.grid.rows, search: true }).render(
      document.getElementById("wrapper")
    );
  },
};
</script>
