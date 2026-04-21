<template>
  <div class="register-container" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div class="register-header">
      <h1>注册</h1>
    </div>
    <form class="register-form" @submit.prevent="register">
      <div class="form-group">
        <label for="username">用户名</label>
        <input type="text" id="username" v-model="form.username" placeholder="请输入用户名" required>
      </div>
      <div class="form-group">
        <label for="password">密码</label>
        <input type="password" id="password" v-model="form.password" placeholder="请输入密码" required>
      </div>
      <div class="form-group">
        <label for="invitationCode">注册码</label>
        <input type="text" id="invitationCode" v-model="form.invitationCode" placeholder="请输入注册码" required>
      </div>
      <button type="submit" class="submit-btn">注册</button>
      <p class="login-link">已有账号？<a @click="goToLogin">立即登录</a></p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useThemeStore } from '../store/theme';
import axios from '../axios';
import { ElMessage } from 'element-plus';

const router = useRouter();
const themeStore = useThemeStore();
const form = ref({
  username: '',
  password: '',
  invitationCode: ''
});

const register = async () => {
  if (!form.value.username || !form.value.password || !form.value.invitationCode) {
    ElMessage.error('请填写所有字段');
    return;
  }

  try {
    const response = await axios.post('/auth/register', {
      username: form.value.username,
      password: form.value.password,
      invitationCode: form.value.invitationCode
    });
    localStorage.setItem('token', response.data.token);
    ElMessage.success('注册成功');
    router.push({ name: 'Blogs' });
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '注册失败');
  }
};

const goToLogin = () => {
  router.push({ name: 'Login' });
};

onMounted(() => {
  themeStore.initTheme();
});
</script>

<style scoped>
.register-container {
  padding: 40px 20px;
  max-width: 800px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.register-container.dark-mode {
  background-color: #000;
  color: #fff;
}

.register-header {
  margin-bottom: 40px;
}

.register-header h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  transition: color 0.3s ease;
}

.register-container.dark-mode .register-header h1 {
  color: #fff;
}

.register-form {
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

.register-container.dark-mode .form-group label {
  color: #fff;
}

.form-group input {
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

.register-container.dark-mode .form-group input {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.form-group input:focus {
  outline: none;
  border-color: #000;
}

.register-container.dark-mode .form-group input:focus {
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

.register-container.dark-mode .submit-btn {
  border-color: #fff;
  background-color: #fff;
  color: #000;
}

.submit-btn:hover {
  background-color: #333;
}

.register-container.dark-mode .submit-btn:hover {
  background-color: #ccc;
}

.login-link {
  text-align: center;
  font-size: 14px;
  color: #888;
}

.register-container.dark-mode .login-link {
  color: #999;
}

.login-link a {
  color: #000;
  cursor: pointer;
  text-decoration: underline;
}

.register-container.dark-mode .login-link a {
  color: #fff;
}

@media (min-width: 768px) {
  .register-container {
    padding: 60px 40px;
  }

  .register-header h1 {
    font-size: 36px;
  }
}
</style>