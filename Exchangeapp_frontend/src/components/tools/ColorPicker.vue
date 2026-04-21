<template>
  <div class="color-picker" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <h2>颜色选择器</h2>
    <div class="picker-content">
      <div class="color-preview">
        <div class="color-display" :style="{ backgroundColor: selectedColor }"></div>
        <input 
          type="color" 
          v-model="selectedColor" 
          class="color-input"
        >
      </div>
      <div class="color-formats">
        <div class="format-item">
          <label>HEX:</label>
          <div class="format-value">
            <input 
              type="text" 
              :value="hexColor" 
              readonly 
              class="format-input"
            >
            <button class="copy-btn" @click="copyColor('hex')">
              复制
            </button>
          </div>
        </div>
        <div class="format-item">
          <label>RGB:</label>
          <div class="format-value">
            <input 
              type="text" 
              :value="rgbColor" 
              readonly 
              class="format-input"
            >
            <button class="copy-btn" @click="copyColor('rgb')">
              复制
            </button>
          </div>
        </div>
        <div class="format-item">
          <label>RGBA:</label>
          <div class="format-value">
            <input 
              type="text" 
              :value="rgbaColor" 
              readonly 
              class="format-input"
            >
            <button class="copy-btn" @click="copyColor('rgba')">
              复制
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useThemeStore } from '../../store/theme';
import { ElMessage } from 'element-plus';

const themeStore = useThemeStore();
const selectedColor = ref('#3498db');

const hexColor = computed(() => selectedColor.value);

const rgbColor = computed(() => {
  const hex = selectedColor.value;
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return `rgb(${r}, ${g}, ${b})`;
});

const rgbaColor = computed(() => {
  const hex = selectedColor.value;
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return `rgba(${r}, ${g}, ${b}, 1)`;
});

const copyColor = async (format: string) => {
  let colorToCopy = '';
  switch (format) {
    case 'hex':
      colorToCopy = hexColor.value;
      break;
    case 'rgb':
      colorToCopy = rgbColor.value;
      break;
    case 'rgba':
      colorToCopy = rgbaColor.value;
      break;
  }
  
  try {
    await navigator.clipboard.writeText(colorToCopy);
    ElMessage.success(`复制 ${format.toUpperCase()} 成功`);
  } catch (error) {
    ElMessage.error('复制失败');
  }
};
</script>

<style scoped>
.color-picker {
  padding: 20px;
  border: 1px solid #000;
  border-radius: 8px;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.color-picker.dark-mode {
  background-color: #000;
  color: #fff;
  border-color: #fff;
}

.color-picker h2 {
  font-size: 20px;
  font-weight: 900;
  margin-bottom: 20px;
  color: #000;
  transition: color 0.3s ease;
}

.color-picker.dark-mode h2 {
  color: #fff;
}

.picker-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.color-preview {
  display: flex;
  align-items: center;
  gap: 20px;
}

.color-display {
  width: 100px;
  height: 100px;
  border: 1px solid #000;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.color-picker.dark-mode .color-display {
  border-color: #fff;
}

.color-input {
  width: 60px;
  height: 40px;
  border: none;
  cursor: pointer;
  background: none;
}

.color-formats {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.format-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

label {
  font-size: 14px;
  font-weight: 700;
  color: #000;
  transition: color 0.3s ease;
}

.color-picker.dark-mode label {
  color: #fff;
}

.format-value {
  display: flex;
  gap: 10px;
}

.format-input {
  flex: 1;
  padding: 8px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 14px;
  font-family: 'Courier New', monospace;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.color-picker.dark-mode .format-input {
  border-color: #fff;
  background-color: #1e1e1e;
  color: #fff;
}

.copy-btn {
  padding: 8px 12px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 14px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
  white-space: nowrap;
}

.color-picker.dark-mode .copy-btn {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.copy-btn:hover {
  background-color: #000;
  color: #fff;
}

.color-picker.dark-mode .copy-btn:hover {
  background-color: #fff;
  color: #000;
}

@media (min-width: 768px) {
  .picker-content {
    flex-direction: row;
    align-items: flex-start;
  }
  
  .color-preview {
    flex-direction: column;
    align-items: center;
  }
  
  .color-formats {
    flex: 1;
  }
}
</style>