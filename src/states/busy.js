import { defineStore } from "pinia";

export const useBusy = defineStore("isloading", {
  state: () => ({
    busy: false,
  }),
  getters: {
    isBusy(state) {
      return state.busy;
    },
  },
  actions: {
    setBusy(busy) {
      this.busy = busy;
    },
  },
});