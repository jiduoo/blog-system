<template>
  <div class="create-blog-container" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div class="create-blog-header">
      <h1>写博客</h1>
    </div>
    <div class="create-blog-form">
      <div class="form-group">
        <label for="title">标题</label>
        <input type="text" id="title" v-model="form.title" placeholder="请输入标题">
      </div>
      <div class="form-group">
        <label for="preview">预览</label>
        <textarea id="preview" v-model="form.preview" placeholder="请输入预览内容" rows="3"></textarea>
      </div>
      
      <!-- Markdown编辑器 -->
      <div class="form-group markdown-editor">
        <div class="editor-header">
          <label for="content">内容 (Markdown)</label>
          <div class="editor-actions">
            <button class="btn" @click="downloadMarkdown" :disabled="!form.content">
              下载 Markdown
            </button>
          </div>
        </div>
        <div class="editor-container">
          <div class="editor-left">
            <textarea 
              id="content" 
              v-model="form.content" 
              placeholder="请输入Markdown内容" 
              rows="15"
            ></textarea>
          </div>
          <div class="editor-right">
            <div class="preview-title">预览</div>
            <div class="markdown-preview" v-html="renderedMarkdown"></div>
          </div>
        </div>
      </div>
      
      <div class="form-group">
        <label for="tags">标签（用逗号分隔）</label>
        <input type="text" id="tags" v-model="tagInput" placeholder="例如：技术,生活,编程">
      </div>
      <div class="form-actions">
        <button class="cancel-btn" @click="goBack">取消</button>
        <button class="submit-btn" @click="submitForm">发布</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from '../axios';
import { useThemeStore } from '../store/theme';
import { ElMessage } from 'element-plus';
import MarkdownIt from 'markdown-it';
import hljs from 'highlight.js';
import 'highlight.js/styles/github.css';

const router = useRouter();
const themeStore = useThemeStore();
const form = ref({
  title: '',
  preview: '',
  content: '',
  tags: [] as string[]
});
const tagInput = ref('');

// 创建markdown-it实例
const md = new MarkdownIt({
  highlight: function(str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value;
      } catch (__) {}
    }
    return ''; // 使用默认的转义
  },
  breaks: true,
  html: true
});

// 实时渲染Markdown
const renderedMarkdown = computed(() => {
  if (!form.value.content) {
    return '<p>请输入Markdown内容</p>';
  }
  return md.render(form.value.content);
});

const submitForm = async () => {
  if (!form.value.title || !form.value.content) {
    ElMessage.error('标题和内容不能为空');
    return;
  }

  // 处理标签输入
  if (tagInput.value) {
    form.value.tags = tagInput.value.split(',').map(tag => tag.trim()).filter(tag => tag);
  }

  try {
    await axios.post('/blogs', form.value);
    ElMessage.success('发布成功');
    router.push({ name: 'Blogs' });
  } catch (error) {
    console.error('发布博客失败:', error);
    ElMessage.error('发布失败，请重试');
  }
};

// 下载Markdown文件
const downloadMarkdown = () => {
  if (!form.value.content) {
    ElMessage.warning('请先输入内容');
    return;
  }

  // 构建Markdown内容
  let markdownContent = `# ${form.value.title}\n\n`;
  if (form.value.preview) {
    markdownContent += `> ${form.value.preview}\n\n`;
  }
  markdownContent += form.value.content;
  if (tagInput.value) {
    const tags = tagInput.value.split(',').map(tag => tag.trim()).filter(tag => tag);
    if (tags.length > 0) {
      markdownContent += `\n\nTags: ${tags.join(', ')}`;
    }
  }

  // 创建Blob对象
  const blob = new Blob([markdownContent], { type: 'text/markdown' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${form.value.title || 'untitled'}.md`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);

  ElMessage.success('Markdown文件已下载');
};

const goBack = () => {
  router.push({ name: 'Blogs' });
};

onMounted(() => {
  themeStore.initTheme();
});
</script>

<style scoped>
.create-blog-container {
  padding: 40px 20px;
  max-width: 1000px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.create-blog-container.dark-mode {
  background-color: #000;
  color: #fff;
}

.create-blog-header {
  margin-bottom: 40px;
}

.create-blog-header h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  transition: color 0.3s ease;
}

.create-blog-container.dark-mode .create-blog-header h1 {
  color: #fff;
}

.create-blog-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 16px;
  font-weight: 700;
  color: #000;
  transition: color 0.3s ease;
}

.create-blog-container.dark-mode .form-group label {
  color: #fff;
}

.form-group input,
.form-group textarea {
  padding: 12px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 16px;
  line-height: 1.5;
  font-family: inherit;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.create-blog-container.dark-mode .form-group input,
.create-blog-container.dark-mode .form-group textarea {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #000;
}

.create-blog-container.dark-mode .form-group input:focus,
.create-blog-container.dark-mode .form-group textarea:focus {
  border-color: #fff;
}

/* Markdown编辑器样式 */
.markdown-editor {
  margin-top: 10px;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.editor-actions {
  display: flex;
  gap: 10px;
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
}

.create-blog-container.dark-mode .btn {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.btn:hover:not(:disabled) {
  background-color: #000;
  color: #fff;
}

.create-blog-container.dark-mode .btn:hover:not(:disabled) {
  background-color: #fff;
  color: #000;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.editor-container {
  display: flex;
  gap: 20px;
  border: 1px solid #000;
  border-radius: 4px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.create-blog-container.dark-mode .editor-container {
  border-color: #fff;
}

.editor-left,
.editor-right {
  flex: 1;
  min-height: 400px;
}

.editor-left textarea {
  width: 100%;
  height: 100%;
  border: none;
  border-right: 1px solid #000;
  border-radius: 0;
  resize: none;
  font-family: 'Courier New', monospace;
}

.create-blog-container.dark-mode .editor-left textarea {
  border-right-color: #fff;
}

.editor-right {
  padding: 20px;
  background-color: #f5f5f5;
  overflow-y: auto;
  transition: all 0.3s ease;
}

.create-blog-container.dark-mode .editor-right {
  background-color: #1e1e1e;
}

.preview-title {
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 16px;
  color: #000;
  transition: color 0.3s ease;
}

.create-blog-container.dark-mode .preview-title {
  color: #fff;
}

.markdown-preview {
  font-size: 16px;
  line-height: 1.8;
  color: #000;
  margin-bottom: 40px;
  transition: color 0.3s ease;
}

.create-blog-container.dark-mode .markdown-preview {
  color: #fff;
}

/* Markdown 样式 */
.markdown-preview h1,
.markdown-preview h2,
.markdown-preview h3,
.markdown-preview h4,
.markdown-preview h5,
.markdown-preview h6 {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
  color: #000;
  transition: color 0.3s ease;
}

.create-blog-container.dark-mode .markdown-preview h1,
.create-blog-container.dark-mode .markdown-preview h2,
.create-blog-container.dark-mode .markdown-preview h3,
.create-blog-container.dark-mode .markdown-preview h4,
.create-blog-container.dark-mode .markdown-preview h5,
.create-blog-container.dark-mode .markdown-preview h6 {
  color: #fff;
}

.markdown-preview h1 {
  font-size: 2em;
}

.markdown-preview h2 {
  font-size: 1.5em;
}

.markdown-preview h3 {
  font-size: 1.25em;
}

.markdown-preview p {
  margin-bottom: 16px;
}

.markdown-preview code {
  background-color: #f6f8fa;
  border-radius: 3px;
  padding: 0.2em 0.4em;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 85%;
  transition: all 0.3s ease;
}

.create-blog-container.dark-mode .markdown-preview code {
  background-color: #2d2d2d;
}

.markdown-preview pre {
  background-color: #f6f8fa;
  border-radius: 6px;
  padding: 16px;
  overflow: auto;
  margin-bottom: 16px;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 85%;
  line-height: 1.45;
  transition: all 0.3s ease;
}

.create-blog-container.dark-mode .markdown-preview pre {
  background-color: #2d2d2d;
}

.markdown-preview pre code {
  background-color: transparent;
  padding: 0;
}

.markdown-preview ul,
.markdown-preview ol {
  margin-left: 24px;
  margin-bottom: 16px;
}

.markdown-preview li {
  margin-bottom: 8px;
}

.markdown-preview a {
  color: #0366d6;
  text-decoration: none;
  transition: color 0.3s ease;
}

.create-blog-container.dark-mode .markdown-preview a {
  color: #64b5f6;
}

.markdown-preview a:hover {
  text-decoration: underline;
}

.markdown-preview blockquote {
  border-left: 4px solid #dfe2e5;
  padding-left: 16px;
  margin: 16px 0;
  color: #6a737d;
  transition: all 0.3s ease;
}

.create-blog-container.dark-mode .markdown-preview blockquote {
  border-left-color: #444;
  color: #999;
}

.markdown-preview img {
  max-width: 100%;
  height: auto;
  margin: 16px 0;
}

/* 代码块容器样式 - Mac风格 */
:deep(.code-block-container) {
  position: relative;
  margin-bottom: 24px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 1px solid #eaeaea;
  background-color: #fff;
  font-family: 'SF Pro Text', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

:deep(.create-blog-container.dark-mode .code-block-container) {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  border-color: #333;
  background-color: #1e1e1e;
}

/* 代码块工具栏 */
:deep(.code-toolbar) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #eaeaea;
  transition: all 0.3s ease;
}

:deep(.create-blog-container.dark-mode .code-toolbar) {
  background-color: #252526;
  border-bottom-color: #333;
}

/* 操作按钮容器 */
:deep(.code-actions) {
  display: flex;
  gap: 6px;
  margin-left: auto;
}

/* 语言标签 */
:deep(.code-lang) {
  font-size: 12px;
  font-weight: 600;
  color: #666;
  transition: color 0.3s ease;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

:deep(.create-blog-container.dark-mode .code-lang) {
  color: #ccc;
}

/* 代码按钮样式 - Mac风格 */
:deep(.code-actions button) {
  background-color: #f0f0f0;
  border: 1px solid #e0e0e0;
  font-size: 10px;
  cursor: pointer;
  padding: 4px 10px;
  border-radius: 4px;
  transition: all 0.2s ease;
  color: #555;
  font-weight: 500;
  font-family: 'SF Pro Text', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  outline: none;
  white-space: nowrap;
  letter-spacing: 0.3px;
  text-transform: uppercase;
}

:deep(.code-actions button:hover) {
  background-color: #e0e0e0;
  border-color: #d0d0d0;
  color: #333;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

:deep(.code-actions button:active) {
  transform: translateY(0);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

/* 深色模式下的按钮样式 */
:deep(.create-blog-container.dark-mode .code-actions button) {
  background-color: #3a3a3a;
  border-color: #4a4a4a;
  color: #ccc;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

:deep(.create-blog-container.dark-mode .code-actions button:hover) {
  background-color: #4a4a4a;
  border-color: #5a5a5a;
  color: #fff;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

:deep(.create-blog-container.dark-mode .code-actions button:active) {
  transform: translateY(0);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

/* 调整代码块样式 */
:deep(.code-block-container pre) {
  margin-bottom: 0;
  border-radius: 0 0 8px 8px;
  border: none;
  margin: 0;
  padding: 16px;
  overflow-x: auto;
  background-color: #282c34;
  color: #abb2bf;
  transition: all 0.3s ease;
}

/* 确保代码块内的代码有正确的字体 */
:deep(.code-block-container pre code) {
  font-family: 'Fira Code', 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', 'Oxygen Mono', 'Ubuntu Monospace', 'Source Code Pro', 'Droid Sans Mono', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
}

/* 深色模式下的代码块样式 */
:deep(.create-blog-container.dark-mode .code-block-container pre) {
  background-color: #1e1e1e;
  color: #d4d4d4;
  border: none;
}

.form-actions {
  display: flex;
  gap: 16px;
  justify-content: flex-end;
  margin-top: 20px;
}

.cancel-btn {
  padding: 12px 24px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 16px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.create-blog-container.dark-mode .cancel-btn {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.cancel-btn:hover {
  border-color: #000;
}

.create-blog-container.dark-mode .cancel-btn:hover {
  border-color: #fff;
}

.submit-btn {
  padding: 12px 24px;
  border: 1px solid #000;
  background-color: #000;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
  cursor: pointer;
  transition: all 0.3s ease;
}

.create-blog-container.dark-mode .submit-btn {
  border-color: #fff;
  background-color: #fff;
  color: #000;
}

.submit-btn:hover {
  background-color: #333;
}

.create-blog-container.dark-mode .submit-btn:hover {
  background-color: #ccc;
}

@media (max-width: 768px) {
  .create-blog-container {
    padding: 20px 16px;
  }
  
  .editor-container {
    flex-direction: column;
  }
  
  .editor-left,
  .editor-right {
    min-height: 300px;
  }
  
  .editor-left textarea {
    border-right: none;
    border-bottom: 1px solid #000;
  }
  
  .create-blog-container.dark-mode .editor-left textarea {
    border-bottom-color: #fff;
  }
}

@media (min-width: 768px) {
  .create-blog-container {
    padding: 60px 40px;
  }
  
  .create-blog-header h1 {
    font-size: 36px;
  }
}
</style>