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

    const curData = context.rootGetters["analytics/data"];
    const mrr = responseData.mrr;
    const months = responseData.months;

    curData.labels = months;

    const newMRR = [];
    const oldMRR = [];
    const expansionMRR = [];
    const reactivationMRR = [];
    const contractionMRR = [];
    const churnMRR = [];

    for (const el of mrr) {
      newMRR.push(el.New);
      oldMRR.push(el.Old);
      expansionMRR.push(el.Expansion);
      reactivationMRR.push(el.Reactivation);
      contractionMRR.push(el.Contraction);
      churnMRR.push(el.Churn);
    }

    for (const el of curData.datasets) {
      if (el.label === "New") {
        el.data = newMRR;
      }
      if (el.label === "Old") {
        el.data = oldMRR;
      }
      if (el.label === "Expansion") {
        el.data = expansionMRR;
      }
      if (el.label === "Reactivation") {
        el.data = reactivationMRR;
      }
      if (el.label === "Contraction") {
        el.data = contractionMRR;
      }
      if (el.label === "Churn") {
        el.data = churnMRR;
      }
    }

    const grid = {
      title: "Monthly Reccuring Revenue (Table)",
      cols: ["", ...months],
      rows: [
        ["New", ...newMRR],
        ["Old", ...oldMRR],
        ["Expansion", ...expansionMRR],
        ["Reactivation", ...reactivationMRR],
        ["Contraction", ...contractionMRR],
        ["Churn", ...churnMRR],
      ],
    };

    context.commit("setGrid", grid);
    context.commit("setData", curData);
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
