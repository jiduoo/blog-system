<template>
  <div class="blog-detail-container" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div class="blog-detail-header">
      <h1>{{ blog.title }}</h1>
      <div class="blog-detail-meta">
        <span>{{ blog.author }}</span>
        <span>{{ formatDate(blog.CreatedAt) }}</span>
        <span>{{ blog.views }} 浏览</span>
      </div>
      <div class="blog-tags" v-if="blog.tags && blog.tags.length > 0">
        <span 
          v-for="tag in blog.tags" 
          :key="tag.ID"
          class="tag"
          @click="filterByTag(tag.Name)"
        >
          {{ tag.Name }}
        </span>
      </div>
    </div>
    <div class="blog-detail-body">
      <!-- 目录栏 -->
      <div class="blog-toc" v-if="toc.length > 0">
        <h3>目录</h3>
        <ul>
          <li 
            v-for="item in toc" 
            :key="item.id"
            :class="'toc-level-' + item.level"
          >
            <a :href="'#' + item.id" @click.prevent="scrollToAnchor(item.id)">
              {{ item.text }}
            </a>
          </li>
        </ul>
      </div>
      <!-- 内容区 -->
      <div class="blog-detail-content" v-html="renderedContent"></div>
    </div>
    <div class="blog-detail-actions">
      <div class="like-btn" @click="likeBlog">
        <span class="like-icon">❤</span>
        <span>{{ blog.likes }} 点赞</span>
      </div>
      <div class="back-btn" @click="goBack">
        返回列表
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, nextTick } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from '../axios';
import { useThemeStore } from '../store/theme';
import MarkdownIt from 'markdown-it';
import hljs from 'highlight.js';
import 'highlight.js/styles/github.css';
import html2canvas from 'html2canvas';

interface Tag {
  ID: number;
  Name: string;
}

interface Blog {
  ID: number;
  title: string;
  author: string;
  CreatedAt: string;
  views: number;
  likes: number;
  content: string;
  tags: Tag[];
}

const router = useRouter();
const route = useRoute();
const themeStore = useThemeStore();
const blog = ref<Blog>({
  ID: 0,
  title: '',
  author: '',
  CreatedAt: '',
  views: 0,
  likes: 0,
  content: '',
  tags: []
});

// 目录数据
const toc = ref<{id: string, text: string, level: number}[]>([]);

// 生成唯一ID
const generateId = (text: string): string => {
  return text.toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, '');
};

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

// 为标题添加ID和目录生成
const defaultHeadingRenderer = md.renderer.rules.heading;
md.renderer.rules.heading = function(tokens, idx, options, env, self) {
  const token = tokens[idx];
  const level = token.tag;
  const text = token.children?.[0]?.content || '';
  const id = generateId(text);
  
  // 添加到目录
  toc.value.push({
    id,
    text,
    level: parseInt(level.replace('h', ''))
  });
  
  // 为标题添加ID
  token.attrSet('id', id);
  
  // 调用默认渲染器
  return defaultHeadingRenderer ? defaultHeadingRenderer(tokens, idx, options, env, self) : self.renderToken(tokens, idx, options);
};

// 渲染Markdown内容
const renderedContent = computed(() => {
  // 每次渲染前清空目录
  toc.value = [];
  
  if (!blog.value.content) {
    return '<p>加载中...</p>';
  }
  return md.render(blog.value.content);
});

// 跳转到锚点
const scrollToAnchor = (id: string) => {
  const element = document.getElementById(id);
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' });
  }
};

const fetchBlogDetail = async () => {
  try {
    const id = route.params.id;
    const response = await axios.get(`/blogs/${id}`);
    blog.value = response.data;
  } catch (error) {
    console.error('获取博客详情失败:', error);
  }
};

const likeBlog = async () => {
  try {
    const id = route.params.id;
    await axios.post(`/blogs/${id}/like`);
    blog.value.likes++;
  } catch (error) {
    console.error('点赞失败:', error);
  }
};

const goBack = () => {
  router.push({ name: 'Blogs' });
};

const filterByTag = async (tagName: string) => {
  router.push({ name: 'Blogs', query: { tag: tagName } });
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString();
};

// 为代码块添加复制、下载等功能
const addCodeBlockActions = () => {
  // 延迟执行，确保DOM完全渲染
  setTimeout(() => {
    const codeBlocks = document.querySelectorAll('.blog-detail-content pre');
    codeBlocks.forEach((codeBlock, index) => {
      // 为代码块添加hljs类，确保代码高亮能够正确应用
      codeBlock.classList.add('hljs');
      
      // 获取代码元素
      const codeElement = codeBlock.querySelector('code');
      if (codeElement) {
        // 为代码元素添加语言类，确保代码高亮能够正确应用
        let lang = 'plaintext';
        const className = codeElement.className;
        const match = className.match(/language-(\w+)/);
        if (match && match[1]) {
          lang = match[1];
        }
        codeElement.classList.add(`language-${lang}`);
        
        // 应用代码高亮
        const highlighted = hljs.highlight(codeElement.textContent || '', { language: lang });
        codeElement.innerHTML = highlighted.value;
      }
      
      // 创建代码块容器
      const container = document.createElement('div');
      container.className = 'code-block-container';
      
      // 移动代码块到容器中
      codeBlock.parentNode?.insertBefore(container, codeBlock);
      container.appendChild(codeBlock);
      
      // 获取语言信息
      let lang = 'plaintext';
      if (codeElement) {
        const className = codeElement.className;
        const match = className.match(/language-(\w+)/);
        if (match && match[1]) {
          lang = match[1];
        }
      }
      
      // 创建工具栏
      const toolbar = document.createElement('div');
      toolbar.className = 'code-toolbar';
      
      // 语言标签
      const langLabel = document.createElement('span');
      langLabel.className = 'code-lang';
      langLabel.textContent = lang;
      toolbar.appendChild(langLabel);
      
      // 操作按钮
      const actions = document.createElement('div');
      actions.className = 'code-actions';
      
      // 复制按钮
      const copyBtn = document.createElement('button');
      copyBtn.innerHTML = '<span>复制</span>';
      copyBtn.title = '复制代码';
      copyBtn.addEventListener('click', () => {
        const code = codeElement?.textContent || '';
        navigator.clipboard.writeText(code).then(() => {
          copyBtn.innerHTML = '<span>已复制</span>';
          setTimeout(() => {
            copyBtn.innerHTML = '<span>复制</span>';
          }, 2000);
        });
      });
      actions.appendChild(copyBtn);
      
      // 下载文件按钮
      const downloadBtn = document.createElement('button');
      downloadBtn.innerHTML = '<span>下载</span>';
      downloadBtn.title = '下载为文件';
      downloadBtn.addEventListener('click', () => {
        const code = codeElement?.textContent || '';
        const blob = new Blob([code], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `code-${index + 1}.${getLanguageExtension(lang)}`;
        a.click();
        URL.revokeObjectURL(url);
      });
      actions.appendChild(downloadBtn);
      
      // 下载截图按钮
      const screenshotBtn = document.createElement('button');
      screenshotBtn.innerHTML = '<span>截图</span>';
      screenshotBtn.title = '下载为截图';
      screenshotBtn.addEventListener('click', async () => {
        try {
          // 禁用按钮，防止重复点击
          screenshotBtn.disabled = true;
          screenshotBtn.innerHTML = '<span>截图中</span>';
          
          // 使用html2canvas截图
          const canvas = await html2canvas(container, {
            backgroundColor: themeStore.isDarkMode ? '#000' : '#fff',
            scale: 2, // 提高截图质量
            logging: false,
            useCORS: true
          });
          
          // 将canvas转换为图片并下载
          const image = canvas.toDataURL('image/png');
          const a = document.createElement('a');
          a.href = image;
          a.download = `code-${index + 1}.png`;
          a.click();
          
          // 恢复按钮状态
          screenshotBtn.innerHTML = '<span>截图</span>';
          screenshotBtn.disabled = false;
        } catch (error) {
          console.error('截图失败:', error);
          screenshotBtn.innerHTML = '<span>截图</span>';
          screenshotBtn.disabled = false;
          alert('截图失败，请重试');
        }
      });
      actions.appendChild(screenshotBtn);
      
      toolbar.appendChild(actions);
      container.insertBefore(toolbar, codeBlock);
    });
  }, 100);
};

// 根据语言获取文件扩展名
const getLanguageExtension = (lang: string): string => {
  const extensions: Record<string, string> = {
    javascript: 'js',
    typescript: 'ts',
    html: 'html',
    css: 'css',
    python: 'py',
    java: 'java',
    c: 'c',
    cpp: 'cpp',
    csharp: 'cs',
    go: 'go',
    rust: 'rs',
    php: 'php',
    ruby: 'rb',
    swift: 'swift',
    kotlin: 'kt',
    json: 'json',
    xml: 'xml',
    yaml: 'yaml',
    markdown: 'md',
    sql: 'sql',
    shell: 'sh',
    powershell: 'ps1',
    batch: 'bat'
  };
  return extensions[lang] || 'txt';
};

onMounted(async () => {
  themeStore.initTheme();
  await fetchBlogDetail();
  await nextTick();
  console.log('Calling addCodeBlockActions');
  addCodeBlockActions();
  console.log('addCodeBlockActions called');
});
</script>

<style>
/* 全局代码高亮样式，确保不受scoped影响 */
pre code {
  display: block;
  overflow-x: auto;
  padding: 1rem;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
}

/* 确保highlight.js样式能够正确应用 */
.hljs {
  background: #f8f9fa !important;
  color: #24292e !important;
  border-radius: 6px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.blog-detail-container.dark-mode .hljs {
  background: #1e1e1e !important;
  color: #d4d4d4 !important;
}

/* 代码高亮主题样式 */
.hljs-keyword, .hljs-selector-tag, .hljs-literal, .hljs-section, .hljs-link {
  color: #0000ff !important;
}

.hljs-string, .hljs-title, .hljs-name, .hljs-type, .hljs-attribute, .hljs-symbol, .hljs-bullet, .hljs-addition, .hljs-variable, .hljs-template-tag, .hljs-template-variable {
  color: #a31515 !important;
}

.hljs-comment, .hljs-quote, .hljs-deletion, .hljs-meta {
  color: #008000 !important;
  font-style: italic;
}

.hljs-number, .hljs-regexp, .hljs-literal, .hljs-built_in, .hljs-builtin-name {
  color: #098658 !important;
}

.hljs-title.class_, .hljs-class .hljs-title {
  color: #2b91af !important;
}

.hljs-attr, .hljs-variable.language_ {
  color: #0000ff !important;
}

.hljs-function .hljs-params {
  color: #000000 !important;
}

/* 深色模式下的代码高亮样式 */
.blog-detail-container.dark-mode .hljs-keyword, 
.blog-detail-container.dark-mode .hljs-selector-tag, 
.blog-detail-container.dark-mode .hljs-literal, 
.blog-detail-container.dark-mode .hljs-section, 
.blog-detail-container.dark-mode .hljs-link {
  color: #569cd6 !important;
}

.blog-detail-container.dark-mode .hljs-string, 
.blog-detail-container.dark-mode .hljs-title, 
.blog-detail-container.dark-mode .hljs-name, 
.blog-detail-container.dark-mode .hljs-type, 
.blog-detail-container.dark-mode .hljs-attribute, 
.blog-detail-container.dark-mode .hljs-symbol, 
.blog-detail-container.dark-mode .hljs-bullet, 
.blog-detail-container.dark-mode .hljs-addition, 
.blog-detail-container.dark-mode .hljs-variable, 
.blog-detail-container.dark-mode .hljs-template-tag, 
.blog-detail-container.dark-mode .hljs-template-variable {
  color: #ce9178 !important;
}

.blog-detail-container.dark-mode .hljs-comment, 
.blog-detail-container.dark-mode .hljs-quote, 
.blog-detail-container.dark-mode .hljs-deletion, 
.blog-detail-container.dark-mode .hljs-meta {
  color: #6a9955 !important;
  font-style: italic;
}

.blog-detail-container.dark-mode .hljs-number, 
.blog-detail-container.dark-mode .hljs-regexp, 
.blog-detail-container.dark-mode .hljs-literal, 
.blog-detail-container.dark-mode .hljs-built_in, 
.blog-detail-container.dark-mode .hljs-builtin-name {
  color: #b5cea8 !important;
}

.blog-detail-container.dark-mode .hljs-title.class_, 
.blog-detail-container.dark-mode .hljs-class .hljs-title {
  color: #4ec9b0 !important;
}

.blog-detail-container.dark-mode .hljs-attr, 
.blog-detail-container.dark-mode .hljs-variable.language_ {
  color: #9cdcfe !important;
}

.blog-detail-container.dark-mode .hljs-function .hljs-params {
  color: #9cdcfe !important;
}
</style>

<style scoped>
.blog-detail-container {
  padding: 40px 20px;
  max-width: 800px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.blog-detail-container.dark-mode {
  background-color: #000;
  color: #fff;
}

.blog-detail-header {
  margin-bottom: 40px;
}

.blog-detail-header h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  margin-bottom: 16px;
  transition: color 0.3s ease;
}

.blog-detail-container.dark-mode .blog-detail-header h1 {
  color: #fff;
}

.blog-detail-meta {
  display: flex;
  gap: 20px;
  font-size: 14px;
  color: #888;
  margin-bottom: 16px;
  transition: color 0.3s ease;
}

.blog-detail-container.dark-mode .blog-detail-meta {
  color: #999;
}

.blog-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 24px;
}

.tag {
  padding: 4px 12px;
  background-color: #f5f5f5;
  border: 1px solid #000;
  border-radius: 16px;
  font-size: 12px;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.blog-detail-container.dark-mode .tag {
  background-color: #1e1e1e;
  border-color: #fff;
  color: #fff;
}

.tag:hover {
  background-color: #333;
  color: #fff;
}

.blog-detail-container.dark-mode .tag:hover {
  background-color: #ccc;
  color: #000;
}

.blog-detail-body {
  display: flex;
  gap: 40px;
  margin-bottom: 40px;
}

/* 目录栏样式 */
.blog-toc {
  width: 250px;
  flex-shrink: 0;
  position: sticky;
  top: 20px;
  max-height: calc(100vh - 100px);
  overflow-y: auto;
  padding: 20px;
  border-radius: 8px;
  background-color: #f8f9fa;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.blog-detail-container.dark-mode .blog-toc {
  background-color: #1e1e1e;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.blog-toc h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #000;
  transition: color 0.3s ease;
}

.blog-detail-container.dark-mode .blog-toc h3 {
  color: #fff;
}

.blog-toc ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.blog-toc li {
  margin-bottom: 8px;
  transition: all 0.3s ease;
}

.blog-toc a {
  color: #666;
  text-decoration: none;
  font-size: 14px;
  transition: all 0.3s ease;
  display: block;
  padding: 4px 0;
}

.blog-detail-container.dark-mode .blog-toc a {
  color: #ccc;
}

.blog-toc a:hover {
  color: #0366d6;
  text-decoration: underline;
}

.blog-detail-container.dark-mode .blog-toc a:hover {
  color: #64b5f6;
}

/* 目录层级缩进 */
.toc-level-2 {
  margin-left: 16px;
}

.toc-level-3 {
  margin-left: 32px;
}

.toc-level-4 {
  margin-left: 48px;
}

.toc-level-5 {
  margin-left: 64px;
}

.toc-level-6 {
  margin-left: 80px;
}

/* 内容区样式 */
.blog-detail-content {
  flex: 1;
  font-size: 16px;
  line-height: 1.8;
  color: #000;
  transition: color 0.3s ease;
}

.blog-detail-container.dark-mode .blog-detail-content {
  color: #fff;
}

/* 响应式布局 */
@media (max-width: 1024px) {
  .blog-detail-body {
    flex-direction: column;
  }
  
  .blog-toc {
    width: 100%;
    position: relative;
    top: 0;
    max-height: none;
    margin-bottom: 24px;
  }
}

/* Markdown 样式 */
.blog-detail-content h1,
.blog-detail-content h2,
.blog-detail-content h3,
.blog-detail-content h4,
.blog-detail-content h5,
.blog-detail-content h6 {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
  color: #000;
  transition: color 0.3s ease;
}

.blog-detail-container.dark-mode .blog-detail-content h1,
.blog-detail-container.dark-mode .blog-detail-content h2,
.blog-detail-container.dark-mode .blog-detail-content h3,
.blog-detail-container.dark-mode .blog-detail-content h4,
.blog-detail-container.dark-mode .blog-detail-content h5,
.blog-detail-container.dark-mode .blog-detail-content h6 {
  color: #fff;
}

.blog-detail-content h1 {
  font-size: 2em;
}

.blog-detail-content h2 {
  font-size: 1.5em;
}

.blog-detail-content h3 {
  font-size: 1.25em;
}

.blog-detail-content p {
  margin-bottom: 16px;
}

.blog-detail-content code {
  background-color: #f6f8fa;
  border-radius: 3px;
  padding: 0.2em 0.4em;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 85%;
  transition: all 0.3s ease;
}

.blog-detail-container.dark-mode .blog-detail-content code {
  background-color: #2d2d2d;
}

.blog-detail-content pre {
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

.blog-detail-container.dark-mode .blog-detail-content pre {
  background-color: #2d2d2d;
}

.blog-detail-content pre code {
  background-color: transparent;
  padding: 0;
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

:deep(.blog-detail-container.dark-mode .code-block-container) {
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

:deep(.blog-detail-container.dark-mode .code-toolbar) {
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

:deep(.blog-detail-container.dark-mode .code-lang) {
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
:deep(.blog-detail-container.dark-mode .code-actions button) {
  background-color: #3a3a3a;
  border-color: #4a4a4a;
  color: #ccc;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

:deep(.blog-detail-container.dark-mode .code-actions button:hover) {
  background-color: #4a4a4a;
  border-color: #5a5a5a;
  color: #fff;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

:deep(.blog-detail-container.dark-mode .code-actions button:active) {
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
:deep(.blog-detail-container.dark-mode .code-block-container pre) {
  background-color: #1e1e1e;
  color: #d4d4d4;
  border: none;
}

.blog-detail-content ul,
.blog-detail-content ol {
  margin-left: 24px;
  margin-bottom: 16px;
}

.blog-detail-content li {
  margin-bottom: 8px;
}

.blog-detail-content a {
  color: #0366d6;
  text-decoration: none;
  transition: color 0.3s ease;
}

.blog-detail-container.dark-mode .blog-detail-content a {
  color: #64b5f6;
}

.blog-detail-content a:hover {
  text-decoration: underline;
}

.blog-detail-content blockquote {
  border-left: 4px solid #dfe2e5;
  padding-left: 16px;
  margin: 16px 0;
  color: #6a737d;
  transition: all 0.3s ease;
}

.blog-detail-container.dark-mode .blog-detail-content blockquote {
  border-left-color: #444;
  color: #999;
}

.blog-detail-content img {
  max-width: 100%;
  height: auto;
  margin: 16px 0;
}

.blog-detail-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20px;
  border-top: 1px solid #000;
  transition: border-color 0.3s ease;
}

.blog-detail-container.dark-mode .blog-detail-actions {
  border-top-color: #fff;
}

.like-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.blog-detail-container.dark-mode .like-btn {
  color: #fff;
}

.like-btn:hover {
  color: #ff4d4f;
}

.like-icon {
  font-size: 20px;
}

.back-btn {
  font-size: 16px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  position: relative;
  transition: color 0.3s ease;
}

.blog-detail-container.dark-mode .back-btn {
  color: #fff;
}

.back-btn::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 100%;
  height: 1px;
  background-color: #000;
  transform: scaleX(0);
  transform-origin: right;
  transition: transform 0.3s ease, background-color 0.3s ease;
}

.blog-detail-container.dark-mode .back-btn::after {
  background-color: #fff;
}

.back-btn:hover::after {
  transform: scaleX(1);
  transform-origin: left;
}

@media (min-width: 768px) {
  .blog-detail-container {
    padding: 60px 40px;
  }
  
  .blog-detail-header h1 {
    font-size: 36px;
  }
}
</style>