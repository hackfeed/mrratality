export default {
  dataIsUploaded(state) {
    return state.dataIsUploaded;
  },
  dataIsLoaded(state) {
    return state.dataIsLoaded;
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
