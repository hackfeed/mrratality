import actions from "./actions.js";
import getters from "./getters.js";
import mutations from "./mutations.js";

export default {
  namespaced: true,
  state() {
    return {
      lastFetch: null,
      files: [],
      file: null,
      data: {
        labels: [
          "5.18",
          "6.18",
          "7.18",
          "8.18",
          "9.18",
          "10.18",
          "11.18",
          "12.18",
          "1.19",
          "2.19",
          "3.19",
          "4.19",
          "5.19",
          "6.19",
        ],
        datasets: [
          {
            label: "New",
            data: [11788.24, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0],
            backgroundColor: "#027436",
          },
          {
            label: "Old",
            data: [
              0.0,
              8623.13,
              17111.85,
              15508.4,
              14079.17,
              15804.03,
              17513.89,
              12551.52,
              11287.53,
              11352.25,
              12768.83,
              12637.74,
              10832.3,
              34.38,
            ],
            backgroundColor: "#09a776",
          },
          {
            label: "Expansion",
            data: [
              0.0,
              17911.03,
              5804.66,
              4289.79,
              5796.42,
              6348.47,
              4210.5,
              2419.02,
              3119.26,
              3959.86,
              4151.3,
              4483.64,
              1317.68,
              0.0,
            ],
            backgroundColor: "#62da9a",
          },
          {
            label: "Reactivation",
            data: [
              0.0,
              0.0,
              496.39,
              963.73,
              1010.91,
              1999.59,
              829.27,
              700.47,
              969.94,
              600.77,
              827.9,
              642.12,
              305.36,
              0.0,
            ],
            backgroundColor: "#707fd7",
          },
          {
            label: "Contraction",
            data: [
              0.0,
              -770.33,
              -5815.2,
              -6001.73,
              -3965.31,
              -3650.42,
              -5151.13,
              -8492.03,
              -3370.68,
              -2654.37,
              -2223.41,
              -4089.72,
              -6175.22,
              -2393.25,
            ],
            backgroundColor: "#ff8700",
          },
          {
            label: "Churn",
            data: [
              0.0,
              -2394.78,
              -3607.11,
              -1902.77,
              -2717.44,
              -1432.06,
              -1487.07,
              -1510.11,
              -1012.79,
              -1370.11,
              -920.64,
              -1020.57,
              -755.97,
              -10027.72,
            ],
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
          text: "Monthly Reccuring Revenue",
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
