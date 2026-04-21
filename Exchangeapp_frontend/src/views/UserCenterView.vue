<template>
  <div class="user-center" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div class="user-center-header">
      <h1>用户中心</h1>
    </div>
    <div class="user-center-content">
      <div class="sidebar">
        <div class="user-info">
          <div class="user-avatar">{{ user?.username?.charAt(0).toUpperCase() }}</div>
          <div class="user-name">{{ user?.username }}</div>
        </div>
        <nav class="nav-menu">
          <button 
            v-for="item in navItems" 
            :key="item.id"
            :class="{ 'active': activeTab === item.id }"
            @click="activeTab = item.id"
            class="nav-item"
          >
            {{ item.name }}
          </button>
        </nav>
      </div>
      <div class="main-content">
        <div v-if="activeTab === 'profile'" class="profile-section">
          <h2>账户信息</h2>
          <div class="form-group">
            <label for="username">用户名</label>
            <input type="text" id="username" v-model="userForm.username" disabled>
          </div>
          <div class="form-group">
            <label for="homePath">个人主页路径</label>
            <div class="homepath-input">
              <span class="prefix">/</span>
              <input 
                type="text" 
                id="homePath" 
                v-model="userForm.homePath"
                placeholder="例如：zhangtest"
              >
            </div>
            <p class="hint">个人主页路径必须唯一，设置后将在网站根目录下访问</p>
          </div>
          <button class="btn" @click="updateProfile" :loading="loading">
            保存设置
          </button>
        </div>
        <div v-if="activeTab === 'password'" class="password-section">
          <h2>修改密码</h2>
          <div class="form-group">
            <label for="oldPassword">旧密码</label>
            <input type="password" id="oldPassword" v-model="passwordForm.oldPassword" placeholder="请输入旧密码">
          </div>
          <div class="form-group">
            <label for="newPassword">新密码</label>
            <input type="password" id="newPassword" v-model="passwordForm.newPassword" placeholder="请输入新密码">
          </div>
          <div class="form-group">
            <label for="confirmPassword">确认新密码</label>
            <input type="password" id="confirmPassword" v-model="passwordForm.confirmPassword" placeholder="请确认新密码">
          </div>
          <button class="btn" @click="updatePassword" :loading="loading">
            修改密码
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { useThemeStore } from '../store/theme';
import { useAuthStore } from '../store/auth';
import axios from '../axios';
import { ElMessage } from 'element-plus';

const themeStore = useThemeStore();
const authStore = useAuthStore();
const activeTab = ref('profile');
const loading = ref(false);
const user = ref<any>(null);

const userForm = reactive({
  username: '',
  homePath: ''
});

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const navItems = [
  { id: 'profile', name: '账户信息' },
  { id: 'password', name: '修改密码' }
];

const fetchUserInfo = async () => {
  try {
    loading.value = true;
    const response = await axios.get('/user/profile');
    user.value = response.data;
    userForm.username = response.data.username;
    userForm.homePath = response.data.homePath || '';
  } catch (error) {
    ElMessage.error('获取用户信息失败');
  } finally {
    loading.value = false;
  }
};

const updateProfile = async () => {
  try {
    if (!userForm.homePath) {
      ElMessage.error('请设置个人主页路径');
      return;
    }
    
    loading.value = true;
    await axios.put('/user/profile', { homePath: userForm.homePath });
    ElMessage.success('个人信息更新成功');
    await fetchUserInfo();
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '更新失败');
  } finally {
    loading.value = false;
  }
};

const updatePassword = async () => {
  try {
    if (!passwordForm.oldPassword || !passwordForm.newPassword || !passwordForm.confirmPassword) {
      ElMessage.error('请填写所有字段');
      return;
    }
    
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      ElMessage.error('两次输入的新密码不一致');
      return;
    }
    
    loading.value = true;
    await axios.put('/user/password', {
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    });
    ElMessage.success('密码修改成功');
    passwordForm.oldPassword = '';
    passwordForm.newPassword = '';
    passwordForm.confirmPassword = '';
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '修改失败');
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  themeStore.initTheme();
  if (authStore.isAuthenticated) {
    fetchUserInfo();
  }
});
</script>

<style scoped>
.user-center {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
  min-height: 80vh;
}

.user-center.dark-mode {
  background-color: #000;
  color: #fff;
}

.user-center-header {
  margin-bottom: 40px;
}

.user-center-header h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  transition: color 0.3s ease;
}

.user-center.dark-mode .user-center-header h1 {
  color: #fff;
}

.user-center-content {
  display: flex;
  gap: 30px;
}

.sidebar {
  width: 250px;
  flex-shrink: 0;
  border: 1px solid #000;
  border-radius: 8px;
  padding: 20px;
  transition: all 0.3s ease;
}

.user-center.dark-mode .sidebar {
  border-color: #fff;
}

.user-info {
  text-align: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #000;
  transition: all 0.3s ease;
}

.user-center.dark-mode .user-info {
  border-bottom-color: #fff;
}

.user-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background-color: #3498db;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 900;
  margin: 0 auto 15px;
}

.user-name {
  font-size: 18px;
  font-weight: 700;
  color: #000;
  transition: color 0.3s ease;
}

.user-center.dark-mode .user-name {
  color: #fff;
}

.nav-menu {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.nav-item {
  padding: 12px 16px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 16px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
  text-align: left;
}

.user-center.dark-mode .nav-item {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.nav-item:hover {
  background-color: #000;
  color: #fff;
}

.user-center.dark-mode .nav-item:hover {
  background-color: #fff;
  color: #000;
}

.nav-item.active {
  background-color: #000;
  color: #fff;
}

.user-center.dark-mode .nav-item.active {
  background-color: #fff;
  color: #000;
}

.main-content {
  flex: 1;
  border: 1px solid #000;
  border-radius: 8px;
  padding: 30px;
  transition: all 0.3s ease;
}

.user-center.dark-mode .main-content {
  border-color: #fff;
}

.profile-section h2,
.password-section h2 {
  font-size: 24px;
  font-weight: 900;
  color: #000;
  margin-bottom: 30px;
  transition: color 0.3s ease;
}

.user-center.dark-mode .profile-section h2,
.user-center.dark-mode .password-section h2 {
  color: #fff;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 16px;
  font-weight: 700;
  color: #000;
  margin-bottom: 8px;
  transition: color 0.3s ease;
}

.user-center.dark-mode .form-group label {
  color: #fff;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 16px;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.user-center.dark-mode .form-group input {
  border-color: #fff;
  background-color: #1e1e1e;
  color: #fff;
}

.homepath-input {
  display: flex;
  align-items: center;
  border: 1px solid #000;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.user-center.dark-mode .homepath-input {
  border-color: #fff;
}

.homepath-input .prefix {
  padding: 0 12px;
  font-size: 16px;
  font-weight: 700;
  color: #000;
  border-right: 1px solid #000;
  transition: all 0.3s ease;
}

.user-center.dark-mode .homepath-input .prefix {
  color: #fff;
  border-right-color: #fff;
}

.homepath-input input {
  border: none;
  flex: 1;
}

.hint {
  font-size: 14px;
  color: #888;
  margin-top: 5px;
  transition: color 0.3s ease;
}

.user-center.dark-mode .hint {
  color: #999;
}

.btn {
  padding: 12px 24px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 16px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
  margin-top: 20px;
}

.user-center.dark-mode .btn {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.btn:hover {
  background-color: #000;
  color: #fff;
}

.user-center.dark-mode .btn:hover {
  background-color: #fff;
  color: #000;
}

@media (max-width: 768px) {
  .user-center-content {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
  }
  
  .nav-menu {
    flex-direction: row;
    overflow-x: auto;
  }
  
  .nav-item {
    white-space: nowrap;
  }
}
</style>