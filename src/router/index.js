import { createRouter, createWebHistory } from "vue-router";
import MainPage from "../components/MainPage.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: MainPage,
    },
    // // Use this style to lazy-load them mfs
    {
      path: "/images",
      name: "images",
      component: () => import("../components/ImagesPage.vue"),
    },
    {
      path: "/session-id",
      name: "session-id",
      component: () => import("../components/SessionId.vue"),
    },
  ],
});

router.beforeEach((to, from, next) => {
  const session_id = localStorage.getItem("session_id");

  if (!session_id && to.path !== "/session-id") {
    // If there is no 'session_id' in local storage and not already on /session-id route
    next("/session-id"); // Redirect to /session-id
  } else {
    next(); // Continue with the navigation
  }
});

export default router;
