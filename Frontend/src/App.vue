<template>
  <div class="page-container" :class="{ 'dark-mode': themeStore.isDarkMode, 'isMobileMenuOpen': isMobileMenuOpen }">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="logo" @click="handleSelect('home')">
        <div class="logo-icon"></div>
        <span>&nbsp;&nbsp;兰佳林的博客</span>
      </div>
      <nav class="nav" v-if="!isMobileMenuOpen">
        <a v-if="authStore.isAuthenticated" @click="handleSelect('userHome')" :class="{ 'active': activeIndex === 'userHome' }">个人主页</a>
        <a @click="handleSelect('blogs')" :class="{ 'active': activeIndex === 'blogs' }">博客</a>
        <a v-if="authStore.isAuthenticated" @click="handleSelect('createBlog')" :class="{ 'active': activeIndex === 'createBlog' }">写博客</a>
        <a @click="handleSelect('navigation')" :class="{ 'active': activeIndex === 'navigation' }">导航</a>
        <a @click="handleSelect('tools')" :class="{ 'active': activeIndex === 'tools' }">工具</a>
        <a v-if="authStore.isAuthenticated" @click="handleSelect('userCenter')" :class="{ 'active': activeIndex === 'UserCenter' }">用户中心</a>
        <a v-if="authStore.isAuthenticated" @click="handleSelect('admin')" :class="{ 'active': activeIndex === 'admin' }">后台管理</a>
        <a v-if="!authStore.isAuthenticated" @click="handleSelect('login')" :class="{ 'active': activeIndex === 'login' }">登录</a>
        <a v-if="!authStore.isAuthenticated" @click="handleSelect('register')" :class="{ 'active': activeIndex === 'register' }">注册</a>
        <a v-if="authStore.isAuthenticated" @click="handleSelect('logout')">退出</a>
        <button class="theme-toggle" @click="themeStore.toggleDarkMode()" aria-label="切换主题">
          <span v-if="!themeStore.isDarkMode">🌙</span>
          <span v-else>☀️</span>
        </button>
      </nav>
      <button class="mobile-menu-btn" @click="toggleMobileMenu" aria-label="菜单">
        <div class="menu-icon"></div>
      </button>
    </header>

    <!-- 移动端菜单 -->
    <div class="mobile-menu" v-if="isMobileMenuOpen">
      <div class="mobile-menu-content">
        <a v-if="authStore.isAuthenticated" @click="handleSelect('userHome')" :class="{ 'active': activeIndex === 'userHome' }">个人主页</a>
        <a @click="handleSelect('blogs')" :class="{ 'active': activeIndex === 'blogs' }">博客</a>
        <a v-if="authStore.isAuthenticated" @click="handleSelect('createBlog')" :class="{ 'active': activeIndex === 'createBlog' }">写博客</a>
        <a @click="handleSelect('navigation')" :class="{ 'active': activeIndex === 'navigation' }">导航</a>
        <a @click="handleSelect('tools')" :class="{ 'active': activeIndex === 'tools' }">工具</a>
        <a v-if="authStore.isAuthenticated" @click="handleSelect('userCenter')" :class="{ 'active': activeIndex === 'UserCenter' }">用户中心</a>
        <a v-if="authStore.isAuthenticated" @click="handleSelect('admin')" :class="{ 'active': activeIndex === 'admin' }">后台管理</a>
        <a v-if="!authStore.isAuthenticated" @click="handleSelect('login')" :class="{ 'active': activeIndex === 'login' }">登录</a>
        <a v-if="!authStore.isAuthenticated" @click="handleSelect('register')" :class="{ 'active': activeIndex === 'register' }">注册</a>
        <a v-if="authStore.isAuthenticated" @click="handleSelect('logout')">退出</a>
        <button class="theme-toggle" @click="themeStore.toggleDarkMode()" aria-label="切换主题">
          <span v-if="!themeStore.isDarkMode">🌙</span>
          <span v-else>☀️</span>
        </button>
      </div>
    </div>

    <!-- 路由视图 → 导航栏下方 -->
    <router-view />

    <!-- 主内容 -->
    <footer class="footer"></footer>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from './store/auth';
import { useThemeStore } from './store/theme';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const themeStore = useThemeStore();
const activeIndex = ref(route.name?.toString() || 'home');
const isMobileMenuOpen = ref(false);
const isMobile = ref(window.innerWidth < 768);

watch(route, (newRoute) => {
  activeIndex.value = newRoute.name?.toString() || 'home';
  isMobileMenuOpen.value = false;
});

const handleSelect = (key: string) => {
  if (key === 'logout') {
    authStore.logout();
    router.push({ name: 'Home' });
  } else if (key === 'userHome') {
    // 跳转到当前用户的个人主页
    // 这里需要先获取用户的 homePath，暂时跳转到一个默认路径
    // 实际使用时，应该从用户信息中获取 homePath
    router.push({ name: 'UserHome', params: { homePath: 'kk3k' } });
  } else {
    router.push({ name: key.charAt(0).toUpperCase() + key.slice(1) });
  }
  isMobileMenuOpen.value = false;
};

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value;
};

const handleResize = () => {
  isMobile.value = window.innerWidth < 768;
  if (!isMobile.value) {
    isMobileMenuOpen.value = false;
  }
};

onMounted(() => {
  themeStore.initTheme();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  transition: all 0.3s ease;
}

/* 全局思源黑体 Heavy = 900 */
.page-container {
  font-family: "Source Han Sans CN", "思源黑体", -apple-system, BlinkMacSystemFont, sans-serif;
  background: #fff;
  color: #000;
  width: 100%;
  min-height: 100vh;
}
.page-container.dark-mode {
  background: #000;
  color: #fff;
}

/* 头部 */
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #000;
  position: relative;
  z-index: 100;
  background-color: #fff;
  transition: background-color 0.3s ease, border-color 0.3s ease;
}
.page-container.dark-mode .header {
  background-color: #000;
  border-bottom-color: #fff;
}
.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 900;
  cursor: pointer;
  transition: color 0.3s ease;
  color: #000;
}
.page-container.dark-mode .logo {
  color: #fff;
}
.logo-icon {
  width: 24px;
  height: 24px;
  position: relative;
  transition: all 0.3s ease;
}
.logo-icon::before,
.logo-icon::after {
  content: "";
  position: absolute;
  width: 16px;
  height: 16px;
  border: 2px solid #000;
  transition: border-color 0.3s ease;
}
.page-container.dark-mode .logo-icon::before,
.page-container.dark-mode .logo-icon::after {
  border-color: #fff;
}
.logo-icon::before {
  top: 0;
  left: 0;
  transform: rotate(45deg);
  border-right: none;
  border-bottom: none;
}
.logo-icon::after {
  bottom: 0;
  right: 0;
  transform: rotate(45deg);
  border-left: none;
  border-top: none;
}
.nav {
  display: flex;
  align-items: center;
  gap: 20px;
}
.nav a {
  text-decoration: none;
  color: #000;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  position: relative;
  transition: all 0.3s ease;
  padding: 8px 0;
}
.page-container.dark-mode .nav a {
  color: #fff;
}
.nav a:hover {
  color: #000;
  transform: translateY(-2px);
}
.page-container.dark-mode .nav a:hover {
  color: #fff;
}
.nav a::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 0;
  height: 2px;
  background-color: #000;
  transition: width 0.3s ease;
}
.page-container.dark-mode .nav a::after {
  background-color: #fff;
}
.nav a:hover::after {
  width: 100%;
}
.nav a.active {
  font-weight: 700;
}
.nav a.active::after {
  width: 100%;
}
.theme-toggle {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}
.theme-toggle:hover {
  background-color: rgba(0, 0, 0, 0.1);
}
.page-container.dark-mode .theme-toggle:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

/* 移动端菜单按钮 */
.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  z-index: 101;
}
.menu-icon {
  width: 24px;
  height: 20px;
  position: relative;
  transition: all 0.3s ease;
}
.menu-icon::before,
.menu-icon::after {
  content: '';
  position: absolute;
  width: 24px;
  height: 2px;
  background-color: #000;
  transition: all 0.3s ease;
}
.page-container.dark-mode .menu-icon::before,
.page-container.dark-mode .menu-icon::after {
  background-color: #fff;
}
.menu-icon::before {
  top: 0;
  left: 0;
}
.menu-icon::after {
  bottom: 0;
  left: 0;
  box-shadow: 0 -8px 0 0 #000;
}
.page-container.dark-mode .menu-icon::after {
  box-shadow: 0 -8px 0 0 #fff;
}
.mobile-menu-btn.active .menu-icon::before {
  transform: rotate(45deg);
  top: 8px;
}
.mobile-menu-btn.active .menu-icon::after {
  transform: rotate(-45deg);
  bottom: 10px;
  box-shadow: none;
}

/* 移动端菜单 */
.mobile-menu {
  position: fixed;
  top: 0;
  right: 0;
  width: 80%;
  height: 100vh;
  background-color: #fff;
  box-shadow: -5px 0 15px rgba(0, 0, 0, 0.1);
  z-index: 99;
  transform: translateX(100%);
  transition: transform 0.3s ease;
}
.page-container.dark-mode .mobile-menu {
  background-color: #000;
  box-shadow: -5px 0 15px rgba(0, 0, 0, 0.3);
}
.mobile-menu-content {
  padding: 80px 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.mobile-menu-content a {
  text-decoration: none;
  color: #000;
  font-size: 18px;
  font-weight: 500;
  padding: 12px 0;
  border-bottom: 1px solid #000;
  transition: all 0.3s ease;
}
.page-container.dark-mode .mobile-menu-content a {
  color: #fff;
  border-bottom-color: #fff;
}
.mobile-menu-content a:hover {
  color: #000;
  transform: translateX(10px);
}
.page-container.dark-mode .mobile-menu-content a:hover {
  color: #fff;
}
.mobile-menu-content a.active {
  font-weight: 700;
  color: #000;
}
.page-container.dark-mode .mobile-menu-content a.active {
  color: #fff;
}
.mobile-menu-content .theme-toggle {
  margin-top: 20px;
  align-self: flex-start;
}

/* 主内容 */


/* 电脑端优化 */
@media (max-width: 767px) {
  .nav {
    display: none;
  }
  .mobile-menu-btn {
    display: block;
  }
  .mobile-menu {
    transform: translateX(100%);
  }
  .page-container.isMobileMenuOpen .mobile-menu {
    transform: translateX(0);
  }
}

@media (min-width: 768px) {
  .header {
    padding: 16px 24px;
  }
  .nav {
    gap: 24px;
  }
  .nav a {
    font-size: 18px;
  }
}

</style>