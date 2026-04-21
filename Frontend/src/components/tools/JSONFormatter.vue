<template>
  <div class="json-formatter" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <h2>JSON 格式化</h2>
    <div class="formatter-content">
      <div class="input-section">
        <label for="json-input">输入 JSON:</label>
        <textarea 
          id="json-input" 
          v-model="inputJson"
          placeholder="请输入 JSON 字符串"
          rows="10"
        ></textarea>
        <div class="input-actions">
          <button class="btn" @click="formatJson">格式化</button>
          <button class="btn" @click="clearInput">清空</button>
        </div>
      </div>
      <div class="output-section">
        <label>格式化结果:</label>
        <pre class="output">{{ formattedJson }}</pre>
        <button class="btn" @click="copyOutput" :disabled="!formattedJson">
          复制结果
        </button>
      </div>
    </div>
    <div class="error-message" v-if="errorMessage">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useThemeStore } from '../../store/theme';
import { ElMessage } from 'element-plus';

const themeStore = useThemeStore();
const inputJson = ref('');
const formattedJson = ref('');
const errorMessage = ref('');

const formatJson = () => {
  try {
    errorMessage.value = '';
    if (!inputJson.value.trim()) {
      errorMessage.value = '请输入 JSON 字符串';
      return;
    }
    const parsed = JSON.parse(inputJson.value);
    formattedJson.value = JSON.stringify(parsed, null, 2);
    ElMessage.success('格式化成功');
  } catch (error) {
    errorMessage.value = 'JSON 格式错误: ' + (error as Error).message;
    ElMessage.error('格式化失败');
  }
};

const clearInput = () => {
  inputJson.value = '';
  formattedJson.value = '';
  errorMessage.value = '';
};

const copyOutput = async () => {
  try {
    await navigator.clipboard.writeText(formattedJson.value);
    ElMessage.success('复制成功');
  } catch (error) {
    ElMessage.error('复制失败');
  }
};
</script>

<style scoped>
.json-formatter {
  padding: 20px;
  border: 1px solid #000;
  border-radius: 8px;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.json-formatter.dark-mode {
  background-color: #000;
  color: #fff;
  border-color: #fff;
}

.json-formatter h2 {
  font-size: 20px;
  font-weight: 900;
  margin-bottom: 20px;
  color: #000;
  transition: color 0.3s ease;
}

.json-formatter.dark-mode h2 {
  color: #fff;
}

.formatter-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-section,
.output-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

label {
  font-size: 14px;
  font-weight: 700;
  color: #000;
  transition: color 0.3s ease;
}

.json-formatter.dark-mode label {
  color: #fff;
}

textarea {
  padding: 12px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 14px;
  font-family: 'Courier New', monospace;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
  resize: vertical;
}

.json-formatter.dark-mode textarea {
  border-color: #fff;
  background-color: #1e1e1e;
  color: #fff;
}

.output {
  padding: 12px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 14px;
  font-family: 'Courier New', monospace;
  background-color: #f5f5f5;
  color: #000;
  transition: all 0.3s ease;
  white-space: pre-wrap;
  max-height: 300px;
  overflow-y: auto;
}

.json-formatter.dark-mode .output {
  border-color: #fff;
  background-color: #1e1e1e;
  color: #fff;
}

.input-actions {
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

.json-formatter.dark-mode .btn {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.btn:hover:not(:disabled) {
  background-color: #000;
  color: #fff;
}

.json-formatter.dark-mode .btn:hover:not(:disabled) {
  background-color: #fff;
  color: #000;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error-message {
  margin-top: 10px;
  padding: 10px;
  background-color: #ffebee;
  color: #c62828;
  border: 1px solid #ef5350;
  border-radius: 4px;
  font-size: 14px;
}

.json-formatter.dark-mode .error-message {
  background-color: #4e342e;
  color: #ffcdd2;
  border-color: #ef5350;
}

@media (min-width: 768px) {
  .formatter-content {
    flex-direction: row;
  }
  
  .input-section,
  .output-section {
    flex: 1;
  }
}
</style>