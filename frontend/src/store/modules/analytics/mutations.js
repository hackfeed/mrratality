export default {
  setFile(state, data) {
    state.file = data;
  },
  setData(state, data) {
    state.data = data;
  },
  setFiles(state, data) {
    state.files = data;
  },
  setFetchTimestamp(state) {
    state.lastFetch = new Date().getTime();
  },
};
