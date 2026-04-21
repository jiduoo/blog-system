<template>
  <div class="base64-tool" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <h2>Base64 编码/解码</h2>
    <div class="tool-content">
      <div class="encode-section">
        <h3>编码</h3>
        <label for="encode-input">输入文本:</label>
        <textarea 
          id="encode-input" 
          v-model="encodeInput"
          placeholder="请输入要编码的文本"
          rows="5"
        ></textarea>
        <button class="btn" @click="encode">编码</button>
      </div>
      <div class="decode-section">
        <h3>解码</h3>
        <label for="decode-input">输入 Base64:</label>
        <textarea 
          id="decode-input" 
          v-model="decodeInput"
          placeholder="请输入要解码的 Base64 字符串"
          rows="5"
        ></textarea>
        <button class="btn" @click="decode">解码</button>
      </div>
      <div class="result-section">
        <h3>结果</h3>
        <textarea 
          v-model="result"
          placeholder="结果将显示在这里"
          rows="5"
          readonly
        ></textarea>
        <button class="btn" @click="copyResult" :disabled="!result">
          复制结果
        </button>
        <button class="btn" @click="clearAll">清空所有</button>
      </div>
    </div>
    <div class="error-message" v-if="errorMessage">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useThemeStore } from '../../store/theme';
import { ElMessage } from 'element-plus';

const themeStore = useThemeStore();
const encodeInput = ref('');
const decodeInput = ref('');
const result = ref('');
const errorMessage = ref('');

const encode = () => {
  try {
    errorMessage.value = '';
    if (!encodeInput.value) {
      errorMessage.value = '请输入要编码的文本';
      return;
    }
    const encoded = btoa(unescape(encodeURIComponent(encodeInput.value)));
    result.value = encoded;
    ElMessage.success('编码成功');
  } catch (error) {
    errorMessage.value = '编码失败: ' + (error as Error).message;
    ElMessage.error('编码失败');
  }
};

const decode = () => {
  try {
    errorMessage.value = '';
    if (!decodeInput.value) {
      errorMessage.value = '请输入要解码的 Base64 字符串';
      return;
    }
    const decoded = decodeURIComponent(escape(atob(decodeInput.value)));
    result.value = decoded;
    ElMessage.success('解码成功');
  } catch (error) {
    errorMessage.value = '解码失败: ' + (error as Error).message;
    ElMessage.error('解码失败');
  }
};

const copyResult = async () => {
  try {
    await navigator.clipboard.writeText(result.value);
    ElMessage.success('复制成功');
  } catch (error) {
    ElMessage.error('复制失败');
  }
};

const clearAll = () => {
  encodeInput.value = '';
  decodeInput.value = '';
  result.value = '';
  errorMessage.value = '';
};
</script>

<style scoped>
.base64-tool {
  padding: 20px;
  border: 1px solid #000;
  border-radius: 8px;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.base64-tool.dark-mode {
  background-color: #000;
  color: #fff;
  border-color: #fff;
}

.base64-tool h2 {
  font-size: 20px;
  font-weight: 900;
  margin-bottom: 20px;
  color: #000;
  transition: color 0.3s ease;
}

.base64-tool.dark-mode h2 {
  color: #fff;
}

.tool-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.encode-section,
.decode-section,
.result-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.encode-section h3,
.decode-section h3,
.result-section h3 {
  font-size: 16px;
  font-weight: 700;
  color: #000;
  transition: color 0.3s ease;
}

.base64-tool.dark-mode .encode-section h3,
.base64-tool.dark-mode .decode-section h3,
.base64-tool.dark-mode .result-section h3 {
  color: #fff;
}

label {
  font-size: 14px;
  font-weight: 700;
  color: #000;
  transition: color 0.3s ease;
}

.base64-tool.dark-mode label {
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

.base64-tool.dark-mode textarea {
  border-color: #fff;
  background-color: #1e1e1e;
  color: #fff;
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
  align-self: flex-start;
}

.base64-tool.dark-mode .btn {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.btn:hover:not(:disabled) {
  background-color: #000;
  color: #fff;
}

.base64-tool.dark-mode .btn:hover:not(:disabled) {
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

.base64-tool.dark-mode .error-message {
  background-color: #4e342e;
  color: #ffcdd2;
  border-color: #ef5350;
}

@media (min-width: 768px) {
  .tool-content {
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto auto;
    grid-template-areas: 
      "encode decode"
      "result result";
  }
  
  .encode-section {
    grid-area: encode;
  }
  
  .decode-section {
    grid-area: decode;
  }
  
  .result-section {
    grid-area: result;
  }
}
</style>