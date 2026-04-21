<template>
  <div class="login-container" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div class="login-header">
      <h1>登录</h1>
    </div>
    <form class="login-form" @submit.prevent="login">
      <div class="form-group">
        <label for="username">用户名</label>
        <input type="text" id="username" v-model="form.username" placeholder="请输入用户名" required>
      </div>
      <div class="form-group">
        <label for="password">密码</label>
        <input type="password" id="password" v-model="form.password" placeholder="请输入密码" required>
      </div>
      <button type="submit" class="submit-btn">登录</button>
      <p class="register-link">还没有账号？<a @click="goToRegister">立即注册</a></p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import { useThemeStore } from '../store/theme';
import { ElMessage } from 'element-plus';

const router = useRouter();
const themeStore = useThemeStore();
const authStore = useAuthStore();
const form = ref({
  username: '',
  password: ''
});

const login = async () => {
  if (!form.value.username || !form.value.password) {
    ElMessage.error('请填写所有字段');
    return;
  }

  try {
    await authStore.login(form.value.username, form.value.password);
    ElMessage.success('登录成功');
    router.push({ name: 'Blogs' });
  } catch (error: any) {
    ElMessage.error('登录失败，请检查用户名和密码');
  }
};

const goToRegister = () => {
  router.push({ name: 'Register' });
};

onMounted(() => {
  themeStore.initTheme();
});
</script>

<style scoped>
.login-container {
  padding: 40px 20px;
  max-width: 800px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.login-container.dark-mode {
  background-color: #000;
  color: #fff;
}

.login-header {
  margin-bottom: 40px;
}

.login-header h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  transition: color 0.3s ease;
}

.login-container.dark-mode .login-header h1 {
  color: #fff;
}

.login-form {
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

.login-container.dark-mode .form-group label {
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

.login-container.dark-mode .form-group input {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.form-group input:focus {
  outline: none;
  border-color: #000;
}

.login-container.dark-mode .form-group input:focus {
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

.login-container.dark-mode .submit-btn {
  border-color: #fff;
  background-color: #fff;
  color: #000;
}

.submit-btn:hover {
  background-color: #333;
}

.login-container.dark-mode .submit-btn:hover {
  background-color: #ccc;
}

.register-link {
  text-align: center;
  font-size: 14px;
  color: #888;
}

.login-container.dark-mode .register-link {
  color: #999;
}

.register-link a {
  color: #000;
  cursor: pointer;
  text-decoration: underline;
}

.login-container.dark-mode .register-link a {
  color: #fff;
}

@media (min-width: 768px) {
  .login-container {
    padding: 60px 40px;
  }

  .login-header h1 {
    font-size: 36px;
  }
}
</style>