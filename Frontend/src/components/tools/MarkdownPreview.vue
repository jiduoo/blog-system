<template>
  <div class="markdown-preview" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <h2>Markdown 预览</h2>
    <div class="preview-content">
      <div class="input-section">
        <label for="markdown-input">输入 Markdown:</label>
        <textarea 
          id="markdown-input" 
          v-model="inputMarkdown"
          placeholder="请输入 Markdown 文本"
          rows="10"
        ></textarea>
        <div class="input-actions">
          <button class="btn" @click="clearInput">清空</button>
        </div>
      </div>
      <div class="output-section">
        <label>预览结果:</label>
        <div class="preview" v-html="renderedMarkdown"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useThemeStore } from '../../store/theme';
import { marked } from 'marked';

const themeStore = useThemeStore();
const inputMarkdown = ref('');

const renderedMarkdown = computed(() => {
  if (!inputMarkdown.value) {
    return '<p>请输入 Markdown 文本</p>';
  }
  return marked(inputMarkdown.value);
});

const clearInput = () => {
  inputMarkdown.value = '';
};
</script>

<style scoped>
.markdown-preview {
  padding: 20px;
  border: 1px solid #000;
  border-radius: 8px;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.markdown-preview.dark-mode {
  background-color: #000;
  color: #fff;
  border-color: #fff;
}

.markdown-preview h2 {
  font-size: 20px;
  font-weight: 900;
  margin-bottom: 20px;
  color: #000;
  transition: color 0.3s ease;
}

.markdown-preview.dark-mode h2 {
  color: #fff;
}

.preview-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-section,
.output-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

label {
  font-size: 14px;
  font-weight: 700;
  color: #000;
  transition: color 0.3s ease;
}

.markdown-preview.dark-mode label {
  color: #fff;
}

textarea {
  padding: 12px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 14px;
  font-family: 'Courier New', monospace;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
  resize: vertical;
}

.markdown-preview.dark-mode textarea {
  border-color: #fff;
  background-color: #1e1e1e;
  color: #fff;
}

.preview {
  padding: 12px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 14px;
  background-color: #f5f5f5;
  color: #000;
  transition: all 0.3s ease;
  min-height: 300px;
  white-space: pre-wrap;
  overflow-y: auto;
}

.markdown-preview.dark-mode .preview {
  border-color: #fff;
  background-color: #1e1e1e;
  color: #fff;
}

.preview :deep(h1),
.preview :deep(h2),
.preview :deep(h3),
.preview :deep(h4),
.preview :deep(h5),
.preview :deep(h6) {
  color: #000;
  margin-top: 20px;
  margin-bottom: 10px;
}

.markdown-preview.dark-mode .preview :deep(h1),
.markdown-preview.dark-mode .preview :deep(h2),
.markdown-preview.dark-mode .preview :deep(h3),
.markdown-preview.dark-mode .preview :deep(h4),
.markdown-preview.dark-mode .preview :deep(h5),
.markdown-preview.dark-mode .preview :deep(h6) {
  color: #fff;
}

.preview :deep(pre) {
  background-color: #f0f0f0;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
}

.markdown-preview.dark-mode .preview :deep(pre) {
  background-color: #2d2d2d;
}

.preview :deep(code) {
  background-color: #f0f0f0;
  padding: 2px 4px;
  border-radius: 2px;
  font-family: 'Courier New', monospace;
}

.markdown-preview.dark-mode .preview :deep(code) {
  background-color: #2d2d2d;
}

.preview :deep(a) {
  color: #1976d2;
  text-decoration: underline;
}

.markdown-preview.dark-mode .preview :deep(a) {
  color: #64b5f6;
}

.preview :deep(ul),
.preview :deep(ol) {
  margin-left: 20px;
  margin-bottom: 10px;
}

.preview :deep(li) {
  margin-bottom: 5px;
}

.btn {
  padding: 8px 16px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 14px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
  align-self: flex-start;
}

.markdown-preview.dark-mode .btn {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.btn:hover {
  background-color: #000;
  color: #fff;
}

.markdown-preview.dark-mode .btn:hover {
  background-color: #fff;
  color: #000;
}

@media (min-width: 768px) {
  .preview-content {
    flex-direction: row;
  }
  
  .input-section,
  .output-section {
    flex: 1;
  }
}
</style>