import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import axios from '../axios';

interface User {
  id: number;
  username: string;
  homePath: string;
  isRoot: boolean;
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'));
  const user = ref<User | null>(null);

  const isAuthenticated = computed(() => !!token.value);

  const fetchUserProfile = async () => {
    if (!token.value) return;
    try {
      const response = await axios.get('/user/profile');
      user.value = response.data;
    } catch (error) {
      console.error('Failed to fetch user profile:', error);
    }
  };

  const login = async (username: string, password: string) => {
    try {
      const response = await axios.post('/auth/login', { username, password });
      token.value = response.data.token;
      localStorage.setItem('token', token.value || '');
      await fetchUserProfile();
    } catch (error) {
      throw new Error(`Login failed! ${error}`);
    }
  };

  const register = async (username: string, password: string, invitationCode: string) => {
    try {
      const response = await axios.post('/auth/register', { username, password, invitationCode });
      token.value = response.data.token;
      localStorage.setItem('token', token.value || '');
      await fetchUserProfile();
    } catch (error) {
      throw new Error(`Register failed! ${error}`);
    }
  };

  const logout = () => {
    token.value = null;
    user.value = null;
    localStorage.removeItem('token');
  };

  return {
    token,
    user,
    isAuthenticated,
    fetchUserProfile,
    login,
    register,
    logout
  };
});