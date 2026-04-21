<template>
  <div class="navigation-container" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div class="navigation-header">
      <h1>导航</h1>
    </div>
    <div class="navigation-content">
      <div class="navigation-section">
        <h2>开发资源</h2>
        <ul class="navigation-list">
          <li><a href="https://developer.mozilla.org/" target="_blank">MDN Web Docs</a></li>
          <li><a href="https://react.dev/" target="_blank">React 文档</a></li>
          <li><a href="https://vuejs.org/" target="_blank">Vue 文档</a></li>
          <li><a href="https://nodejs.org/" target="_blank">Node.js 文档</a></li>
        </ul>
      </div>
      <div class="navigation-section">
        <h2>设计资源</h2>
        <ul class="navigation-list">
          <li><a href="https://dribbble.com/" target="_blank">Dribbble</a></li>
          <li><a href="https://behance.net/" target="_blank">Behance</a></li>
          <li><a href="https://coolors.co/" target="_blank">Coolors</a></li>
          <li><a href="https://unsplash.com/" target="_blank">Unsplash</a></li>
        </ul>
      </div>
      <div class="navigation-section">
        <h2>工具网站</h2>
        <ul class="navigation-list">
          <li><a href="https://codepen.io/" target="_blank">CodePen</a></li>
          <li><a href="https://stackblitz.com/" target="_blank">StackBlitz</a></li>
          <li><a href="https://codesandbox.io/" target="_blank">CodeSandbox</a></li>
          <li><a href="https://regex101.com/" target="_blank">RegExr</a></li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useThemeStore } from '../store/theme';

const themeStore = useThemeStore();

onMounted(() => {
  themeStore.initTheme();
});
</script>

<style scoped>
.navigation-container {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.navigation-container.dark-mode {
  background-color: #000;
  color: #fff;
}

.navigation-header {
  margin-bottom: 40px;
}

.navigation-header h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  transition: color 0.3s ease;
}

.navigation-container.dark-mode .navigation-header h1 {
  color: #fff;
}

.navigation-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: 40px;
}

.navigation-section h2 {
  font-size: 20px;
  font-weight: 900;
  color: #000;
  margin-bottom: 16px;
  transition: color 0.3s ease;
}

.navigation-container.dark-mode .navigation-section h2 {
  color: #fff;
}

.navigation-list {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.navigation-list li a {
  text-decoration: none;
  color: #000;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.3s ease;
  position: relative;
  padding-left: 16px;
}

.navigation-container.dark-mode .navigation-list li a {
  color: #fff;
}

.navigation-list li a::before {
  content: '→';
  position: absolute;
  left: 0;
  transition: transform 0.3s ease;
}

.navigation-list li a:hover {
  color: #000;
}

.navigation-container.dark-mode .navigation-list li a:hover {
  color: #fff;
}

.navigation-list li a:hover::before {
  transform: translateX(4px);
}

@media (min-width: 768px) {
  .navigation-container {
    padding: 60px 40px;
  }

  .navigation-header h1 {
    font-size: 36px;
  }

  .navigation-content {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>