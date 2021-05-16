export default {
  async uploadData(context, data) {
    const response = await fetch("http://localhost:8081/files/upload", {
      method: "POST",
      body: data,
      headers: {
        token: localStorage.getItem("token"),
      },
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to fill data");
      throw error;
    }

    context.commit("setFile", responseData.filename);
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
  },
  async loadFiles(context) {
    const response = await fetch("http://localhost:8081/files/load", {
      headers: {
        token: localStorage.getItem("token"),
      },
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to load files");
      throw error;
    }

    context.commit("setFiles", responseData.files);
  },
};
