<template>  
    <main class="main" :class="{ 'dark-mode': themeStore.isDarkMode }">
          <div class="hero-text">
            <div class="greeting-zh">你好，</div>
            <div class="greeting-en">hello</div>

            <!-- 重点：这一排永远不换行，横向滚动 -->
            <div class="categories">
              <div class="category-item" v-for="category in categories" :key="category">{{ category }}</div>
            </div>

            <div class="world-text">world</div>
            <div class="dot-wrapper">
                <div class="dot">。</div>
            </div>
          </div>

          <div class="read-btn" @click="goToBlogs">
            <div class="read-btn-text">开启阅读</div>
            <div class="read-btn-arrow"></div>
          </div>
        </main>
        <hr>
</template>  

<script setup lang="ts">  
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useThemeStore } from '../store/theme';

const router = useRouter();
const themeStore = useThemeStore();
const categories = ref([
  '世界', '前端', '后端', '运维', '安全',
  '动画', '漫画', '音乐', '插画', '生活'
]);

const goToBlogs = () => {
  router.push({ name: 'Blogs' });
};

onMounted(() => {
  themeStore.initTheme();
});
</script>  

<style scoped>  
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  transition: all 0.3s ease;
}
.main {
  display: flex;
  flex-direction: column;
  gap: 40px;
  padding: 40px 20px 60px;
  min-height: calc(100vh - 56px);
  background-color: #fff;
  color: #222;
}
.main.dark-mode {
  background-color: #121212;
  color: #e0e0e0;
}
.hero-text {
  width: 100%;
}
.greeting-zh {
  font-size: 72px;
  line-height: 1.1;
  font-weight: 900;
  color: #000;
  writing-mode: vertical-lr;
  transition: color 0.3s ease;
}
.main.dark-mode .greeting-zh {
  color: #fff;
}

.greeting-en {
  font-size: 28px;
  color: #888;
  margin: 8px 0 16px;
  font-weight: 700;
  transition: color 0.3s ease;
}
.main.dark-mode .greeting-en {
  color: #999;
}

/* —————————— 核心：永远一排、不换行 —————————— */
.categories {
  margin-top: 30px;
  display: flex;
  flex-direction: row;     /* 横向 */
  flex-wrap: nowrap;      /* 禁止换行 */
  gap: 16px;              /* 间距 */
  overflow-x: hidden;        
  /* 超出滚动 */
  padding-bottom: 8px;
  width: 100%;
  white-space: nowrap;
}
.category-item {
  font-size: 22px;
  font-weight: 900;
  flex-shrink: 0; 
  writing-mode: vertical-lr;        /* 禁止压缩，保证不换行 */
  cursor: pointer;
  position: relative;
  transition: transform 0.3s ease, color 0.3s ease;
  color: #000;
}
.main.dark-mode .category-item {
  color: #fff;
}
.category-item:hover {
  transform: translateY(-5px);
  color: #000;
}
.main.dark-mode .category-item:hover {
  color: #fff;
}
.category-item::after {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 2px;
  background-color: #000;
  transition: width 0.3s ease;
}
.main.dark-mode .category-item::after {
  background-color: #fff;
}
.category-item:hover::after {
  width: 80%;
}

.world-text {
  font-size: 28px;
  color: #888;
  margin-top: 20px;
  font-weight: 700;
  transition: color 0.3s ease;
}
.main.dark-mode .world-text {
  color: #999;
}
.dot-wrapper {
  width: 30%;
  text-align: right;
  padding-right: 10px; /* 可选，微调边距更美观 */
}
.dot {
  font-family: "Source Han Sans CN", "思源黑体", -apple-system, BlinkMacSystemFont, sans-serif;
  display: inline-block;
  width: 16px;
  height: 16px;
  font-size: 72px;
  line-height: 1.1;
  font-weight: 900;
  color: #000;
  /* border-radius: 50%;
  background: #333; */
  margin-top: 20px;
  transition: color 0.3s ease;
}
.main.dark-mode .dot {
  color: #fff;
}

/* 阅读按钮 */
.read-btn {
  border: 1px solid #000;
  padding: 16px 20px;
  width: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  background-color: #fff;
}
.main.dark-mode .read-btn {
  border-color: #fff;
  background-color: #000;
}
.read-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  background-color: #f5f5f5;
}
.main.dark-mode .read-btn:hover {
  box-shadow: 0 5px 15px rgba(255, 255, 255, 0.1);
  background-color: #333;
}
.read-btn-text {
  font-size: 22px;
  font-weight: 900;
  transition: color 0.3s ease;
  color: #000;
}
.main.dark-mode .read-btn-text {
  color: #fff;
}
.read-btn-arrow {
  width: 20px;
  height: 20px;
  background:
    linear-gradient(45deg, transparent 45%, #000 45%, #000 55%, transparent 55%),
    linear-gradient(135deg, transparent 45%, #000 45%, #000 55%, transparent 55%);
  transition: background 0.3s ease;
}
.main.dark-mode .read-btn-arrow {
  background:
    linear-gradient(45deg, transparent 45%, #fff 45%, #fff 55%, transparent 55%),
    linear-gradient(135deg, transparent 45%, #fff 45%, #fff 55%, transparent 55%);
}
.read-btn:hover .read-btn-arrow {
  transform: translateX(5px);
  transition: transform 0.3s ease;
}

.footer {
  height: 2px;
  background: #000;
  margin: 0 16px 20px;
  transition: background 0.3s ease;
}
.main.dark-mode .footer {
  background: #fff;
}

@media (min-width: 768px) {
  .main {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    padding: 100px 12%;
  }
  .read-btn {
    flex-direction: column;
    width: auto;
    padding: 64px 16px;
  }
  .read-btn-text {
    writing-mode: vertical-rl;
    margin-bottom: 24px;
  }
  .dot-wrapper {
    width: 60%;
  }
}

</style>