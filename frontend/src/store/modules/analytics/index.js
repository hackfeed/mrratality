import actions from "./actions.js";
import getters from "./getters.js";
import mutations from "./mutations.js";

export default {
  namespaced: true,
  state() {
    return {
      files: [],
      file: null,
      grid: {
        title: null,
        cols: null,
        rows: null,
      },
      data: {
        labels: [],
        datasets: [
          {
            label: "New",
            data: [],
            backgroundColor: "#027436",
          },
          {
            label: "Old",
            data: [],
            backgroundColor: "#09a776",
          },
          {
            label: "Expansion",
            data: [],
            backgroundColor: "#62da9a",
          },
          {
            label: "Reactivation",
            data: [],
            backgroundColor: "#707fd7",
          },
          {
            label: "Contraction",
            data: [],
            backgroundColor: "#ff8700",
          },
          {
            label: "Churn",
            data: [],
            backgroundColor: "#8f0239",
          },
        ],
      },
      dataOptions: {
        responsive: true,
        legend: {
          display: false,
        },
        title: {
          display: true,
          text: "Monthly Reccuring Revenue (Chart)",
          fontSize: 24,
          fontColor: "black",
        },
        tooltips: {
          backgroundColor: "#17BF62",
        },
        scales: {
          xAxes: [
            {
              stacked: true,
              gridLines: {
                display: false,
              },
            },
          ],
          yAxes: [
            {
              stacked: true,
              ticks: {
                beginAtZero: true,
              },
              gridLines: {
                display: false,
              },
            },
          ],
        },
      },
    };
  },
  actions,
  getters,
  mutations,
};
