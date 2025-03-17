import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router"
import HomeView from "../views/HomeView.vue"
import StockDetailView from "../views/StockDetailView.vue"
import RecommendationsView from "../views/RecommendationsView.vue"

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "home",
    component: HomeView,
  },
  {
    path: "/stock/:ticker",
    name: "stock-detail",
    component: StockDetailView,
    props: true,
  },
  {
    path: "/recommendations",
    name: "recommendations",
    component: RecommendationsView,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router

