import { defineStore } from 'pinia';
import { ref, onMounted } from 'vue';

export const useThemeStore = defineStore('theme', () => {
  const isDarkMode = ref(false);

  const toggleDarkMode = () => {
    isDarkMode.value = !isDarkMode.value;
    updateTheme();
  };

  const updateTheme = () => {
    if (isDarkMode.value) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
    localStorage.setItem('darkMode', isDarkMode.value.toString());
  };

  const initTheme = () => {
    const savedDarkMode = localStorage.getItem('darkMode');
    if (savedDarkMode) {
      isDarkMode.value = savedDarkMode === 'true';
      updateTheme();
    }
  };

  onMounted(() => {
    initTheme();
  });

  return {
    isDarkMode,
    toggleDarkMode,
    initTheme
  };
});
