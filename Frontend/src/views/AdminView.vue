<template>
  <div class="admin-container" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div class="admin-header">
      <h1>后台管理</h1>
    </div>
    <div class="admin-content">
      <div class="admin-nav">
        <button 
          :class="{ 'active': activeTab === 'blogs' }"
          @click="activeTab = 'blogs'"
        >
          文章管理
        </button>
      </div>
      
      <!-- 文章管理 -->
      <div v-if="activeTab === 'blogs'" class="admin-section">
        <h2>文章列表</h2>
        
        <!-- 搜索功能 -->
        <div class="search-form">
          <input 
            v-model="searchKeyword"
            placeholder="搜索文章标题、内容或作者"
            @keyup.enter="searchBlogs"
          >
          <button @click="searchBlogs">搜索</button>
          <button @click="clearSearch">清除</button>
        </div>
        
        <div class="blog-list">
          <div 
            v-for="blog in blogs" 
            :key="blog.ID"
            class="blog-item"
          >
            <div class="blog-header">
              <h3>{{ blog.title }}</h3>
              <div class="blog-actions">
                <button @click="editBlog(blog)">编辑</button>
                <button class="delete-btn" @click="deleteBlog(blog.ID)">删除</button>
              </div>
            </div>
            <p class="blog-meta">
              作者: {{ blog.author }} | 创建时间: {{ formatDate(blog.CreatedAt) }} | 浏览: {{ blog.views }} | 点赞: {{ blog.likes }}
            </p>
            <div class="blog-tags">
              <span 
                v-for="tag in blog.tags" 
                :key="tag.ID"
                class="tag"
              >
                {{ tag.Name }}
              </span>
              <span v-if="!blog.tags || blog.tags.length === 0" class="no-tags">无标签</span>
            </div>
            <p class="blog-preview">{{ blog.preview }}</p>
          </div>
        </div>
        
        <div v-if="blogs.length === 0" class="empty-message">
          暂无文章
        </div>
      </div>
    </div>
    
    <!-- 编辑博客弹窗 -->
    <div v-if="showEditModal" class="modal-overlay" @click.self="closeEditModal">
      <div class="modal-content">
        <h2>编辑文章</h2>
        <div class="form-group">
          <label>标题</label>
          <input v-model="editingBlog.title" type="text">
        </div>
        <div class="form-group">
          <label>预览</label>
          <textarea v-model="editingBlog.preview" rows="3"></textarea>
        </div>
        <div class="form-group">
          <label>内容 (Markdown)</label>
          <textarea v-model="editingBlog.content" rows="10"></textarea>
        </div>
        <div class="form-group">
          <label>标签（用逗号分隔）</label>
          <input v-model="editingBlog.tagsInput" type="text" placeholder="例如：技术,生活,编程">
        </div>
        <div class="modal-actions">
          <button @click="closeEditModal">取消</button>
          <button class="submit-btn" @click="saveBlog">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useThemeStore } from '../store/theme';
import axios from '../axios';
import { ElMessage, ElMessageBox } from 'element-plus';

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
  preview: string;
  content: string;
  tags: Tag[];
}

const themeStore = useThemeStore();
const activeTab = ref('blogs');
const blogs = ref<Blog[]>([]);
const searchKeyword = ref('');

// 编辑博客相关
const showEditModal = ref(false);
const editingBlog = ref({
  ID: 0,
  title: '',
  preview: '',
  content: '',
  tagsInput: ''
});

const loadBlogs = async () => {
  try {
    const response = await axios.get('/blogs');
    blogs.value = response.data;
  } catch (error) {
    ElMessage.error('加载文章失败');
  }
};

const searchBlogs = async () => {
  try {
    if (searchKeyword.value.trim()) {
      const response = await axios.get(`/blogs/search?keyword=${encodeURIComponent(searchKeyword.value.trim())}`);
      blogs.value = response.data;
    } else {
      loadBlogs();
    }
  } catch (error) {
    ElMessage.error('搜索文章失败');
  }
};

const clearSearch = () => {
  searchKeyword.value = '';
  loadBlogs();
};

const editBlog = (blog: Blog) => {
  editingBlog.value = {
    ID: blog.ID,
    title: blog.title,
    preview: blog.preview,
    content: blog.content,
    tagsInput: blog.tags?.map(t => t.Name).join(', ') || ''
  };
  showEditModal.value = true;
};

const saveBlog = async () => {
  if (!editingBlog.value.title.trim() || !editingBlog.value.content.trim()) {
    ElMessage.warning('标题和内容不能为空');
    return;
  }

  try {
    const tagsArray = editingBlog.value.tagsInput
      .split(',')
      .map(t => t.trim())
      .filter(t => t);

    await axios.put(`/blogs/${editingBlog.value.ID}`, {
      title: editingBlog.value.title,
      preview: editingBlog.value.preview,
      content: editingBlog.value.content,
      tags: tagsArray
    });
    
    ElMessage.success('文章更新成功');
    closeEditModal();
    loadBlogs();
  } catch (error) {
    ElMessage.error('更新文章失败');
  }
};

const deleteBlog = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await axios.delete(`/blogs/${id}`);
    ElMessage.success('文章删除成功');
    loadBlogs();
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除文章失败');
    }
  }
};

const closeEditModal = () => {
  showEditModal.value = false;
  editingBlog.value = {
    ID: 0,
    title: '',
    preview: '',
    content: '',
    tagsInput: ''
  };
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleString();
};

onMounted(() => {
  themeStore.initTheme();
  loadBlogs();
});
</script>

<style scoped>
.admin-container {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.admin-container.dark-mode {
  background-color: #000;
  color: #fff;
}

.admin-header {
  margin-bottom: 40px;
}

.admin-header h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  transition: color 0.3s ease;
}

.admin-container.dark-mode .admin-header h1 {
  color: #fff;
}

.admin-content {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.admin-nav {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.admin-nav button {
  padding: 12px 24px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 16px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .admin-nav button {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.admin-nav button:hover {
  background-color: #333;
  color: #fff;
}

.admin-container.dark-mode .admin-nav button:hover {
  background-color: #ccc;
  color: #000;
}

.admin-nav button.active {
  background-color: #000;
  color: #fff;
}

.admin-container.dark-mode .admin-nav button.active {
  background-color: #fff;
  color: #000;
}

.admin-section h2 {
  font-size: 24px;
  font-weight: 900;
  color: #000;
  margin-bottom: 24px;
  transition: color 0.3s ease;
}

.admin-container.dark-mode .admin-section h2 {
  color: #fff;
}

/* 搜索功能 */
.search-form {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.search-form input {
  flex: 1;
  padding: 12px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 16px;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .search-form input {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.search-form button {
  padding: 12px 24px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 14px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .search-form button {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.search-form button:hover {
  background-color: #333;
  color: #fff;
}

.admin-container.dark-mode .search-form button:hover {
  background-color: #ccc;
  color: #000;
}

.blog-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.blog-item {
  padding: 24px;
  border: 1px solid #000;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .blog-item {
  border-color: #fff;
}

.blog-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.blog-item h3 {
  font-size: 20px;
  font-weight: 700;
  color: #000;
  margin-bottom: 8px;
  transition: color 0.3s ease;
}

.admin-container.dark-mode .blog-item h3 {
  color: #fff;
}

.blog-actions {
  display: flex;
  gap: 8px;
}

.blog-actions button {
  padding: 6px 12px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 12px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .blog-actions button {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.blog-actions button:hover {
  background-color: #333;
  color: #fff;
}

.admin-container.dark-mode .blog-actions button:hover {
  background-color: #ccc;
  color: #000;
}

.blog-actions .delete-btn:hover {
  background-color: #ff4d4f;
  border-color: #ff4d4f;
  color: #fff;
}

.admin-container.dark-mode .blog-actions .delete-btn:hover {
  background-color: #ff4d4f;
  border-color: #ff4d4f;
  color: #fff;
}

.blog-meta {
  font-size: 14px;
  color: #888;
  margin-bottom: 12px;
  transition: color 0.3s ease;
}

.admin-container.dark-mode .blog-meta {
  color: #999;
}

.blog-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.tag {
  padding: 4px 12px;
  background-color: #f5f5f5;
  border: 1px solid #000;
  border-radius: 16px;
  font-size: 12px;
  color: #000;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .tag {
  background-color: #1e1e1e;
  border-color: #fff;
  color: #fff;
}

.no-tags {
  font-size: 12px;
  color: #888;
  font-style: italic;
}

.admin-container.dark-mode .no-tags {
  color: #999;
}

.blog-preview {
  font-size: 14px;
  line-height: 1.6;
  color: #000;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.3s ease;
}

.admin-container.dark-mode .blog-preview {
  color: #fff;
}

.empty-message {
  text-align: center;
  padding: 40px;
  font-size: 16px;
  color: #888;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: #fff;
  padding: 32px;
  border-radius: 8px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  transition: all 0.3s ease;
}

.dark-mode .modal-content {
  background-color: #000;
  border: 1px solid #fff;
}

.modal-content h2 {
  font-size: 24px;
  font-weight: 900;
  color: #000;
  margin-bottom: 24px;
}

.dark-mode .modal-content h2 {
  color: #fff;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 700;
  color: #000;
  margin-bottom: 8px;
}

.dark-mode .form-group label {
  color: #fff;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 14px;
  font-family: inherit;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.dark-mode .form-group input,
.dark-mode .form-group textarea {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.form-group textarea {
  resize: vertical;
  font-family: 'Courier New', monospace;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.modal-actions button {
  padding: 12px 24px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 14px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.dark-mode .modal-actions button {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.modal-actions button:hover {
  background-color: #333;
  color: #fff;
}

.dark-mode .modal-actions button:hover {
  background-color: #ccc;
  color: #000;
}

.modal-actions .submit-btn {
  background-color: #000;
  color: #fff;
}

.dark-mode .modal-actions .submit-btn {
  background-color: #fff;
  color: #000;
}

.modal-actions .submit-btn:hover {
  background-color: #333;
  color: #fff;
}

.dark-mode .modal-actions .submit-btn:hover {
  background-color: #ccc;
}

@media (min-width: 768px) {
  .admin-container {
    padding: 60px 40px;
  }

  .admin-header h1 {
    font-size: 36px;
  }
}
</style>
