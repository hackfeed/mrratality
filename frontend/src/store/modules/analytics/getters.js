export default {
  files(state) {
    return state.files;
  },
  file(state) {
    return state.file;
  },
  data(state) {
    return state.data;
  },
  dataOptions(state) {
    return state.dataOptions;
  },
  shouldUpdate(state) {
    const lastFetch = state.lastFetch;
    if (!lastFetch) {
      return true;
    }
    const currentTimestamp = new Date().getTime();
    return (currentTimestamp - lastFetch) / 1000 > 60;
  },
};
