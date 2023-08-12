import { defineStore } from "pinia";

export const useSession = defineStore("sessionId", {
  state: () => ({
    session: "",
  }),
  getters: {
    getSession(state) {
      state.session = localStorage.getItem("session_id");
      return state.session;
    },
  },
  actions: {
    setSession(session) {
      this.session = session;
    },
  },
});