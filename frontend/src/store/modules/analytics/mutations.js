export default {
  setDataIsUploaded(state, data) {
    state.dataIsUploaded = data;
  },
  setDataIsLoaded(state, data) {
    state.dataIsLoaded = data;
  },
  setData(state, data) {
    state.data = data;
  },
  setFetchTimestamp(state) {
    state.lastFetch = new Date().getTime();
  },
};
