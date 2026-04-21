<template>
  <div class="tools-container" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div class="tools-header">
      <h1>工具</h1>
    </div>
    <div class="tools-nav">
      <button 
        v-for="tool in tools" 
        :key="tool.id"
        :class="{ 'active': activeTool === tool.id }"
        @click="activeTool = tool.id"
        class="nav-btn"
      >
        {{ tool.name }}
      </button>
    </div>
    <div class="tool-content">
      <JSONFormatter v-if="activeTool === 'json'" />
      <MarkdownPreview v-else-if="activeTool === 'markdown'" />
      <ColorPicker v-else-if="activeTool === 'color'" />
      <Base64Tool v-else-if="activeTool === 'base64'" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useThemeStore } from '../store/theme';
import JSONFormatter from '../components/tools/JSONFormatter.vue';
import MarkdownPreview from '../components/tools/MarkdownPreview.vue';
import ColorPicker from '../components/tools/ColorPicker.vue';
import Base64Tool from '../components/tools/Base64Tool.vue';

const themeStore = useThemeStore();
const activeTool = ref('json');

const tools = [
  { id: 'json', name: 'JSON 格式化' },
  { id: 'markdown', name: 'Markdown 预览' },
  { id: 'color', name: '颜色选择器' },
  { id: 'base64', name: 'Base64 编码/解码' }
];

onMounted(() => {
  themeStore.initTheme();
});
</script>

<style scoped>
.tools-container {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.tools-container.dark-mode {
  background-color: #000;
  color: #fff;
}

.tools-header {
  margin-bottom: 40px;
}

.tools-header h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  transition: color 0.3s ease;
}

.tools-container.dark-mode .tools-header h1 {
  color: #fff;
}

.tools-nav {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 30px;
  border-bottom: 1px solid #000;
  padding-bottom: 10px;
  transition: all 0.3s ease;
}

.tools-container.dark-mode .tools-nav {
  border-bottom-color: #fff;
}

.nav-btn {
  padding: 10px 20px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 14px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.tools-container.dark-mode .nav-btn {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.nav-btn:hover {
  background-color: #000;
  color: #fff;
}

.tools-container.dark-mode .nav-btn:hover {
  background-color: #fff;
  color: #000;
}

.nav-btn.active {
  background-color: #000;
  color: #fff;
}

.tools-container.dark-mode .nav-btn.active {
  background-color: #fff;
  color: #000;
}

.tool-content {
  min-height: 500px;
}

@media (min-width: 768px) {
  .tools-container {
    padding: 60px 40px;
  }

  .tools-header h1 {
    font-size: 36px;
  }
}
</style>