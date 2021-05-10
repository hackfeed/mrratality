export default {
  async uploadData(context, data) {
    const response = await fetch("http://localhost:8081/upload", {
      method: "POST",
      body: data,
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to fill data.");
      throw error;
    }

    context.commit("setDataIsUploaded", responseData.isUploaded);
  },
  async loadData(context) {
    // const periodStart = data.periodStart;
    // const periodEnd = data.periodEnd;

    // const response = await fetch("path/to/backend", {
    //   body: JSON.stringify({ periodStart, periodEnd }),
    // });
    // const responseData = await response.json();

    // if (!response.ok) {
    //   const error = new Error(responseData.message || "Failed to fill data.");
    //   throw error;
    // }

    context.commit("setData", context.rootGetters["analytics/data"]);
    context.commit("setDataIsLoaded", true);
  },
};
