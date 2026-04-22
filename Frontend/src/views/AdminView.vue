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
        <button 
          v-if="isRoot"
          :class="{ 'active': activeTab === 'users' }"
          @click="activeTab = 'users'"
        >
          用户管理
        </button>
        <button 
          v-if="isRoot"
          :class="{ 'active': activeTab === 'invitations' }"
          @click="activeTab = 'invitations'"
        >
          注册码管理
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
      
      <!-- 用户管理 -->
      <div v-if="activeTab === 'users' && isRoot" class="admin-section">
        <h2>用户管理</h2>
        
        <div class="admin-actions">
          <button class="add-btn" @click="showCreateUserModal = true">添加用户</button>
        </div>
        
        <div class="user-list">
          <div 
            v-for="user in users" 
            :key="user.ID"
            class="user-item"
          >
            <div class="user-info">
              <h3>{{ user.Username }}</h3>
              <p class="user-meta">
                ID: {{ user.ID }} | {{ user.IsRoot ? '管理员' : '普通用户' }}
              </p>
            </div>
            <div class="user-actions">
              <button @click="editUser(user)">编辑</button>
              <button 
                v-if="!user.IsRoot"
                class="delete-btn" 
                @click="deleteUser(user.ID)"
              >
                删除
              </button>
            </div>
          </div>
        </div>
        
        <div v-if="users.length === 0" class="empty-message">
          暂无用户
        </div>
      </div>
      
      <!-- 注册码管理 -->
      <div v-if="activeTab === 'invitations' && isRoot" class="admin-section">
        <h2>注册码管理</h2>
        
        <div class="admin-actions">
          <button class="add-btn generate-btn" @click="generateInvitationCode">生成注册码</button>
          <button class="cleanup-btn" @click="cleanupExpiredCodes">清理过期码</button>
        </div>
        
        <div class="invitation-list">
          <div 
            v-for="code in invitationCodes" 
            :key="code.ID"
            class="invitation-item"
            :class="{ 'expired': isExpired(code.ExpiresAt), 'used': code.Used }"
          >
            <div class="invitation-info">
              <h3>{{ code.Code }}</h3>
              <p class="invitation-meta">
                创建者: {{ code.CreatedBy }} | 创建时间: {{ formatDate(code.CreatedAt) }}
              </p>
              <p class="invitation-meta">
                过期时间: {{ formatDate(code.ExpiresAt) }} | 状态: {{ code.Used ? '已使用' : '未使用' }}
              </p>
              <p v-if="code.Used" class="invitation-meta">
                使用人: {{ code.UsedBy }}
              </p>
            </div>
            <div class="invitation-actions">
              <button class="delete-btn" @click="deleteInvitationCode(code.Code)">删除</button>
            </div>
          </div>
        </div>
        
        <div v-if="invitationCodes.length === 0" class="empty-message">
          暂无注册码
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
    
    <!-- 创建用户弹窗 -->
    <div v-if="showCreateUserModal" class="modal-overlay" @click.self="closeCreateUserModal">
      <div class="modal-content">
        <h2>添加用户</h2>
        <div class="form-group">
          <label>用户名</label>
          <input v-model="newUser.username" type="text" placeholder="请输入用户名">
        </div>
        <div class="form-group">
          <label>密码</label>
          <input v-model="newUser.password" type="password" placeholder="请输入密码">
        </div>
        <div class="form-group">
          <label>是否为管理员</label>
          <input v-model="newUser.isRoot" type="checkbox">
        </div>
        <div class="modal-actions">
          <button @click="closeCreateUserModal">取消</button>
          <button class="submit-btn" @click="createUser">保存</button>
        </div>
      </div>
    </div>
    
    <!-- 编辑用户弹窗 -->
    <div v-if="showEditUserModal" class="modal-overlay" @click.self="closeEditUserModal">
      <div class="modal-content">
        <h2>编辑用户</h2>
        <div class="form-group">
          <label>用户名</label>
          <input v-model="editingUser.username" type="text" placeholder="请输入用户名">
        </div>
        <div class="form-group">
          <label>新密码（留空表示不修改）</label>
          <input v-model="editingUser.password" type="password" placeholder="请输入新密码">
        </div>
        <div class="form-group">
          <label>是否为管理员</label>
          <input v-model="editingUser.isRoot" type="checkbox" :disabled="editingUser.id === 1">
          <span v-if="editingUser.id === 1" class="disabled-hint">(root用户不可修改)</span>
        </div>
        <div class="modal-actions">
          <button @click="closeEditUserModal">取消</button>
          <button class="submit-btn" @click="updateUser">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useThemeStore } from '../store/theme';
import { useAuthStore } from '../store/auth';
import axios from '../axios';
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus';

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

interface User {
  ID: number;
  Username: string;
  IsRoot: boolean;
  HomePath: string;
}

interface InvitationCode {
  ID: number;
  Code: string;
  Used: boolean;
  CreatedBy: string;
  CreatedAt: string;
  ExpiresAt: string;
  UsedBy: string;
}

const themeStore = useThemeStore();
const authStore = useAuthStore();
const activeTab = ref('blogs');
const blogs = ref<Blog[]>([]);
const users = ref<User[]>([]);
const invitationCodes = ref<InvitationCode[]>([]);
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

// 编辑用户相关
const showCreateUserModal = ref(false);
const showEditUserModal = ref(false);
const newUser = ref({
  username: '',
  password: '',
  isRoot: false
});
const editingUser = ref({
  id: 0,
  username: '',
  password: '',
  isRoot: false
});

// 计算属性
const isRoot = computed(() => authStore.user?.isRoot || false);

// 加载文章
const loadBlogs = async () => {
  try {
    const response = await axios.get('/blogs');
    blogs.value = response.data;
  } catch (error) {
    ElMessage.error('加载文章失败');
  }
};

// 加载用户
const loadUsers = async () => {
  try {
    const response = await axios.get('/users');
    users.value = response.data;
  } catch (error) {
    ElMessage.error('加载用户失败');
  }
};

// 加载邀请码
const loadInvitationCodes = async () => {
  try {
    const response = await axios.get('/invitation-codes');
    invitationCodes.value = response.data;
  } catch (error) {
    ElMessage.error('加载注册码失败');
  }
};

// 搜索文章
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

// 清除搜索
const clearSearch = () => {
  searchKeyword.value = '';
  loadBlogs();
};

// 编辑文章
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

// 保存文章
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

// 删除文章
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

// 关闭编辑文章弹窗
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

// 关闭创建用户弹窗
const closeCreateUserModal = () => {
  showCreateUserModal.value = false;
  newUser.value = {
    username: '',
    password: '',
    isRoot: false
  };
};

// 关闭编辑用户弹窗
const closeEditUserModal = () => {
  showEditUserModal.value = false;
  editingUser.value = {
    id: 0,
    username: '',
    password: '',
    isRoot: false
  };
};

// 编辑用户
const editUser = (user: User) => {
  editingUser.value = {
    id: user.ID,
    username: user.Username,
    password: '',
    isRoot: user.IsRoot
  };
  showEditUserModal.value = true;
};

// 创建用户
const createUser = async () => {
  if (!newUser.value.username.trim() || !newUser.value.password.trim()) {
    ElMessage.warning('用户名和密码不能为空');
    return;
  }

  try {
    await axios.post('/users', {
      username: newUser.value.username,
      password: newUser.value.password,
      isRoot: newUser.value.isRoot
    });
    
    ElMessage.success('用户创建成功');
    closeCreateUserModal();
    loadUsers();
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '创建用户失败');
  }
};

// 更新用户
const updateUser = async () => {
  if (!editingUser.value.username.trim()) {
    ElMessage.warning('用户名不能为空');
    return;
  }

  try {
    await axios.put(`/users/${editingUser.value.id}`, {
      username: editingUser.value.username,
      password: editingUser.value.password,
      isRoot: editingUser.value.isRoot
    });
    
    ElMessage.success('用户更新成功');
    closeEditUserModal();
    loadUsers();
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '更新用户失败');
  }
};

// 删除用户
const deleteUser = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个用户吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await axios.delete(`/users/${id}`);
    ElMessage.success('用户删除成功');
    loadUsers();
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除用户失败');
    }
  }
};

// 生成注册码
const generateInvitationCode = async () => {
  try {
    const response = await axios.post('/generate-invitation-code');
    
    ElNotification({
      title: '成功',
      message: `生成的注册码: ${response.data.code}`,
      type: 'success',
      duration: 5000
    });
    loadInvitationCodes();
  } catch (error) {
    ElMessage.error('生成注册码失败');
  }
};

// 删除注册码
const deleteInvitationCode = async (code: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个注册码吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await axios.delete(`/invitation-codes/${code}`);
    ElMessage.success('注册码删除成功');
    loadInvitationCodes();
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除注册码失败');
    }
  }
};

// 清理过期注册码
const cleanupExpiredCodes = async () => {
  try {
    await axios.delete('/invitation-codes/cleanup');
    ElMessage.success('过期注册码清理成功');
    loadInvitationCodes();
  } catch (error) {
    ElMessage.error('清理过期注册码失败');
  }
};

// 检查注册码是否过期
const isExpired = (expiresAt: string) => {
  return new Date(expiresAt) < new Date();
};

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleString();
};

// 切换标签时加载对应数据
const handleTabChange = (tab: string) => {
  activeTab.value = tab;
  if (tab === 'users' && isRoot.value) {
    loadUsers();
  } else if (tab === 'invitations' && isRoot.value) {
    loadInvitationCodes();
  }
};

onMounted(async () => {
  themeStore.initTheme();
  await authStore.fetchUserProfile();
  console.log('Is root in onMounted after fetch:', isRoot.value);
  loadBlogs();
  if (isRoot.value) {
    loadUsers();
    loadInvitationCodes();
  }
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

/* 管理操作按钮 */
.admin-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.admin-actions button {
  padding: 12px 24px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 16px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .admin-actions button {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.admin-actions button:hover {
  background-color: #333;
  color: #fff;
}

.admin-container.dark-mode .admin-actions button:hover {
  background-color: #ccc;
  color: #000;
}

.admin-actions .add-btn {
  background-color: #4CAF50;
  border-color: #4CAF50;
  color: white;
}

.admin-actions .add-btn:hover {
  background-color: #45a049;
  color: white;
}

.admin-container.dark-mode .admin-actions .add-btn {
  background-color: #4CAF50;
  border-color: #4CAF50;
  color: white;
}

.admin-container.dark-mode .admin-actions .add-btn:hover {
  background-color: #45a049;
  color: white;
}

.admin-actions .cleanup-btn {
  background-color: #ff9800;
  border-color: #ff9800;
  color: white;
}

.admin-actions .cleanup-btn:hover {
  background-color: #f57c00;
  color: white;
}

.admin-container.dark-mode .admin-actions .cleanup-btn {
  background-color: #ff9800;
  border-color: #ff9800;
  color: white;
}

.admin-container.dark-mode .admin-actions .cleanup-btn:hover {
  background-color: #f57c00;
  color: white;
}

.admin-actions .generate-btn {
  background-color: #2196F3;
  border-color: #2196F3;
  color: white;
}

.admin-actions .generate-btn:hover {
  background-color: #1976D2;
  color: white;
}

.admin-container.dark-mode .admin-actions .generate-btn {
  background-color: #2196F3;
  border-color: #2196F3;
  color: white;
}

.admin-container.dark-mode .admin-actions .generate-btn:hover {
  background-color: #1976D2;
  color: white;
}

/* 用户列表 */
.user-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.user-item {
  padding: 20px;
  border: 1px solid #000;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .user-item {
  border-color: #fff;
}

.user-info h3 {
  font-size: 18px;
  font-weight: 700;
  color: #000;
  margin-bottom: 8px;
  transition: color 0.3s ease;
}

.admin-container.dark-mode .user-info h3 {
  color: #fff;
}

.user-meta {
  font-size: 14px;
  color: #888;
  transition: color 0.3s ease;
}

.admin-container.dark-mode .user-meta {
  color: #999;
}

.user-actions {
  display: flex;
  gap: 8px;
}

.user-actions button {
  padding: 6px 12px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 12px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .user-actions button {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.user-actions button:hover {
  background-color: #333;
  color: #fff;
}

.admin-container.dark-mode .user-actions button:hover {
  background-color: #ccc;
  color: #000;
}

.user-actions .delete-btn:hover {
  background-color: #ff4d4f;
  border-color: #ff4d4f;
  color: #fff;
}

.admin-container.dark-mode .user-actions .delete-btn:hover {
  background-color: #ff4d4f;
  border-color: #ff4d4f;
  color: #fff;
}

/* 邀请码列表 */
.invitation-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.invitation-item {
  padding: 20px;
  border: 1px solid #000;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .invitation-item {
  border-color: #fff;
}

.invitation-item.expired {
  opacity: 0.6;
  border-color: #ff9800;
}

.invitation-item.used {
  opacity: 0.6;
  border-color: #9e9e9e;
}

.invitation-info h3 {
  font-size: 18px;
  font-weight: 700;
  color: #000;
  margin-bottom: 8px;
  word-break: break-all;
  transition: color 0.3s ease;
}

.admin-container.dark-mode .invitation-info h3 {
  color: #fff;
}

.invitation-meta {
  font-size: 14px;
  color: #888;
  margin-bottom: 4px;
  transition: color 0.3s ease;
}

.admin-container.dark-mode .invitation-meta {
  color: #999;
}

.invitation-actions {
  display: flex;
  gap: 8px;
}

.invitation-actions button {
  padding: 6px 12px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 12px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-container.dark-mode .invitation-actions button {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.invitation-actions button:hover {
  background-color: #333;
  color: #fff;
}

.admin-container.dark-mode .invitation-actions button:hover {
  background-color: #ccc;
  color: #000;
}

.invitation-actions .delete-btn:hover {
  background-color: #ff4d4f;
  border-color: #ff4d4f;
  color: #fff;
}

.admin-container.dark-mode .invitation-actions .delete-btn:hover {
  background-color: #ff4d4f;
  border-color: #ff4d4f;
  color: #fff;
}

/* 禁用提示 */
.disabled-hint {
  font-size: 12px;
  color: #888;
  margin-left: 8px;
  font-style: italic;
}

.admin-container.dark-mode .disabled-hint {
  color: #999;
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
