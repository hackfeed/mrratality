<template>
  <div>
    <section class="centered">
      <base-card>
        <h2>Get full control of your app's money flow with mrratality</h2>
        <base-button link to="/analytics">Try now</base-button>
      </base-card>
      <base-card>
        <h2>Quick MRR report</h2>
        <img src="@/assets/report.png" alt="Quick MRR Report" />
      </base-card>
      <base-card>
        <h2>Flexible period adjustment</h2>
        <img src="@/assets/periods.png" alt="Flexible period adjustment" />
      </base-card>
      <base-card>
        <h2>Simple data export from CSV</h2>
        <img src="@/assets/import.png" alt="Simple data export from CSV" />
      </base-card>
      <base-card>
        <h2>Upload your data now and get carte blanche of your money flow</h2>
        <base-button link to="/analytics">Upload data</base-button>
      </base-card>
    </section>
  </div>
</template>

<script>
export default {};
</script>

<style scoped>
img {
  height: 20rem;
}
.centered {
  text-align: center;
}
</style>
