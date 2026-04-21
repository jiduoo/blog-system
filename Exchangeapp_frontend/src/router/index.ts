import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import Login from '../components/Login.vue';
import Register from '../components/Register.vue';
import BlogsView from '../views/BlogsView.vue';
import BlogDetailView from '../views/BlogDetailView.vue';
import CreateBlogView from '../views/CreateBlogView.vue';
import NavigationView from '../views/NavigationView.vue';
import ToolsView from '../views/ToolsView.vue';
import AdminView from '../views/AdminView.vue';
import UserCenterView from '../views/UserCenterView.vue';
import UserHomeView from '../views/UserHomeView.vue';

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'Home', component: HomeView },
  { path: '/blogs', name: 'Blogs', component: BlogsView },
  { path: '/blogs/:id', name: 'BlogDetail', component: BlogDetailView },
  { path: '/blog/:id', redirect: to => `/blogs/${to.params.id}` },
  { path: '/create-blog', name: 'CreateBlog', component: CreateBlogView },
  { path: '/navigation', name: 'Navigation', component: NavigationView },
  { path: '/tools', name: 'Tools', component: ToolsView },
  { path: '/login', name: 'Login', component: Login },
  { path: '/register', name: 'Register', component: Register },
  { path: '/admin', name: 'Admin', component: AdminView },
  { path: '/user-center', name: 'UserCenter', component: UserCenterView },
  { path: '/:homePath', name: 'UserHome', component: UserHomeView },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
