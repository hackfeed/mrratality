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
  async loadData(context, data) {
    const response = await fetch("http://localhost:8081/analytics", {
      method: "POST",
      body: JSON.stringify(data),
      headers: {
        token: localStorage.getItem("token"),
      },
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to load analytics.");
      throw error;
    }

    console.log(responseData);

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
  async deleteFile(context, data) {
    const reqData = {
      name: data,
    };

    const response = await fetch("http://localhost:8081/files/delete", {
      method: "POST",
      body: JSON.stringify(reqData),
      headers: {
        token: localStorage.getItem("token"),
      },
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to delete file");
      throw error;
    }

    const files = context.rootGetters["analytics/files"].filter(
      (file) => file.name != reqData.name
    );

    context.commit("setFiles", files);
  },
};
