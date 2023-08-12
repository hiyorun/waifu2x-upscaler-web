import { defineStore } from "pinia";

export const useBusy = defineStore("isloading", {
  state: () => ({
    busy:false,
    assetsToLoad:0,
    assetsLoaded:0,
  }),
  getters: {
    loadingProgress(state) {
      return (state.assetsToLoad / state.assetsLoaded) * 100 || 0;
    },
    isBusy(state){
      return state.busy
    }
  },
  actions: {
    loadAssets(count,done){
      this.assetsLoaded = done
      this.assetsToLoad = count
    },
    resetCounter() {
      this.assetsLoaded = 0
      this.assetsToLoad = 0
    },
    setBusy(busy){
      this.busy = busy
    }
  },
});

// ⠁⠂⠄⡀⢀⠠⠐⠈