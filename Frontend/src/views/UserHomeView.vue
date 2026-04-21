<template>
  <div class="user-home" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div v-if="user" class="user-home-content">
      <div class="user-profile">
        <div class="user-avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
        <h1>{{ user.username }}的主页</h1>
        <p class="user-description">欢迎访问{{ user.username }}的个人主页</p>
      </div>
      <div class="user-info">
        <h2>个人信息</h2>
        <div class="info-item">
          <span class="label">用户名:</span>
          <span class="value">{{ user.username }}</span>
        </div>
        <div class="info-item">
          <span class="label">个人主页:</span>
          <span class="value">{{ origin }}/{{ user.homePath }}</span>
        </div>
      </div>
      
      <div class="user-blogs">
        <div class="blogs-header">
          <h2>{{ user.username }}的博客</h2>
          <div class="search-box">
            <input 
              type="text" 
              v-model="searchKeyword" 
              placeholder="搜索博客..."
              @keyup.enter="searchBlogs"
            />
            <button @click="searchBlogs">搜索</button>
          </div>
        </div>
        
        <div class="blogs-content">
          <!-- 左侧时间列表 -->
          <div class="time-list">
            <h3>时间归档</h3>
            <div class="time-filter">
              <div 
                class="time-item" 
                :class="{ 'active': !selectedDate }"
                @click="clearDateFilter"
              >
                全部时间
              </div>
            </div>
            <div class="year-list">
              <div 
                v-for="year in years" 
                :key="year.year"
                class="year-item"
              >
                <div 
                  class="year-header" 
                  @click="toggleYear(year.year)"
                >
                  <span>{{ year.year }}</span>
                  <span class="expand-icon">{{ expandedYears.includes(year.year) ? '▼' : '▶' }}</span>
                </div>
                <div 
                  v-if="expandedYears.includes(year.year)"
                  class="date-list"
                >
                  <div 
                    v-for="date in year.dates" 
                    :key="date.date"
                    class="date-item"
                    :class="{ 'active': selectedDate === date.date }"
                    @click="filterByDate(date.date)"
                  >
                    {{ date.date }} ({{ date.count }})
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 右侧博客内容 -->
          <div class="main-content">
            <div class="blogs-list">
              <div v-for="blog in filteredBlogs" :key="blog.ID" class="blog-item">
                <div class="blog-title">{{ blog.title }}</div>
                <div class="blog-meta">
                  <span>{{ formatDate(blog.CreatedAt) }}</span>
                  <span>{{ blog.views }} 浏览</span>
                  <span>{{ blog.likes }} 点赞</span>
                </div>
                <div class="blog-tags" v-if="blog.tags && blog.tags.length > 0">
                  <span 
                    v-for="tag in blog.tags" 
                    :key="tag.ID"
                    class="tag"
                  >
                    {{ tag.Name }}
                  </span>
                </div>
                <div class="blog-preview">{{ blog.preview }}</div>
                <div class="blog-read-more" @click="navigateToDetail(blog.ID)">阅读更多</div>
              </div>
            </div>
            
            <div v-if="filteredBlogs.length === 0" class="empty-message">
              {{ selectedDate ? '该日期没有文章' : searchKeyword ? '未找到匹配的文章' : '暂无博客' }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="not-found">
      <h1>用户不存在</h1>
      <p>该个人主页路径不存在或已被删除</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useThemeStore } from '../store/theme';
import axios from '../axios';

interface DateGroup {
  date: string;
  count: number;
}

interface YearGroup {
  year: string;
  dates: DateGroup[];
}

const themeStore = useThemeStore();
const route = useRoute();
const router = useRouter();
const user = ref<any>(null);
const blogs = ref<any[]>([]);
const origin = window.location.origin;
const searchKeyword = ref('');
const selectedDate = ref('');
const expandedYears = ref<string[]>([]);

const fetchUserByHomePath = async () => {
  try {
    const homePath = route.params.homePath as string;
    const response = await axios.get(`/user/home/${homePath}`);
    user.value = response.data;
    if (user.value) {
      await fetchUserBlogs(user.value.id);
    }
  } catch (error) {
    user.value = null;
  }
};

const fetchUserBlogs = async (userId: number) => {
  try {
    // 假设后端有一个接口可以根据用户ID获取博客
    // 这里暂时使用通用的博客列表接口，后续可以添加专门的接口
    const response = await axios.get('/blogs');
    // 过滤出当前用户的博客
    blogs.value = response.data.filter((blog: any) => blog.author === user.value.username);
  } catch (error) {
    console.error('获取用户博客失败:', error);
    blogs.value = [];
  }
};

const searchBlogs = async () => {
  if (searchKeyword.value) {
    try {
      const response = await axios.get(`/blogs/search?keyword=${encodeURIComponent(searchKeyword.value)}`);
      // 过滤出当前用户的博客
      blogs.value = response.data.filter((blog: any) => blog.author === user.value.username);
    } catch (error) {
      console.error('搜索博客失败:', error);
    }
  } else {
    await fetchUserBlogs(user.value.id);
  }
  selectedDate.value = '';
};

const filterByDate = (date: string) => {
  selectedDate.value = date;
  searchKeyword.value = '';
};

const clearDateFilter = async () => {
  selectedDate.value = '';
  if (searchKeyword.value) {
    await searchBlogs();
  } else {
    await fetchUserBlogs(user.value.id);
  }
};

const toggleYear = (year: string) => {
  const index = expandedYears.value.indexOf(year);
  if (index > -1) {
    expandedYears.value.splice(index, 1);
  } else {
    expandedYears.value.push(year);
  }
};

const years = computed(() => {
  const yearMap = new Map<string, Set<string>>();
  
  blogs.value.forEach(blog => {
    const date = new Date(blog.CreatedAt);
    const year = date.getFullYear().toString();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    const dateStr = `${month}-${day}`;
    
    if (!yearMap.has(year)) {
      yearMap.set(year, new Set());
    }
    yearMap.get(year)?.add(dateStr);
  });
  
  const result: YearGroup[] = [];
  yearMap.forEach((dates, year) => {
    const dateGroups: DateGroup[] = [];
    dates.forEach(date => {
      const count = blogs.value.filter(blog => {
        const blogDate = new Date(blog.CreatedAt);
        const blogMonth = (blogDate.getMonth() + 1).toString().padStart(2, '0');
        const blogDay = blogDate.getDate().toString().padStart(2, '0');
        return blogMonth + '-' + blogDay === date;
      }).length;
      dateGroups.push({ date, count });
    });
    
    // 按日期排序
    dateGroups.sort((a, b) => a.date.localeCompare(b.date));
    
    result.push({ year, dates: dateGroups });
  });
  
  // 按年份降序排序
  result.sort((a, b) => parseInt(b.year) - parseInt(a.year));
  
  return result;
});

const filteredBlogs = computed(() => {
  if (selectedDate.value) {
    return blogs.value.filter(blog => {
      const blogDate = new Date(blog.CreatedAt);
      const month = (blogDate.getMonth() + 1).toString().padStart(2, '0');
      const day = blogDate.getDate().toString().padStart(2, '0');
      const dateStr = `${month}-${day}`;
      return dateStr === selectedDate.value;
    });
  }
  return blogs.value;
});

const navigateToDetail = (id: number) => {
  router.push({ name: 'BlogDetail', params: { id } });
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString();
};

onMounted(() => {
  themeStore.initTheme();
  fetchUserByHomePath();
});
</script>

<style scoped>
.user-home {
  padding: 40px 20px;
  max-width: 800px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
  min-height: 80vh;
}

.user-home.dark-mode {
  background-color: #000;
  color: #fff;
}

.user-profile {
  text-align: center;
  margin-bottom: 40px;
  padding-bottom: 30px;
  border-bottom: 1px solid #000;
  transition: all 0.3s ease;
}

.user-home.dark-mode .user-profile {
  border-bottom-color: #fff;
}

.user-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background-color: #3498db;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: 900;
  margin: 0 auto 20px;
}

.user-profile h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  margin-bottom: 10px;
  transition: color 0.3s ease;
}

.user-home.dark-mode .user-profile h1 {
  color: #fff;
}

.user-description {
  font-size: 18px;
  color: #888;
  transition: color 0.3s ease;
}

.user-home.dark-mode .user-description {
  color: #999;
}

.user-info {
  margin-top: 30px;
}

.user-info h2 {
  font-size: 24px;
  font-weight: 900;
  color: #000;
  margin-bottom: 20px;
  transition: color 0.3s ease;
}

.user-home.dark-mode .user-info h2 {
  color: #fff;
}

.info-item {
  display: flex;
  margin-bottom: 15px;
  padding: 12px;
  border: 1px solid #000;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.user-home.dark-mode .info-item {
  border-color: #fff;
}

.label {
  width: 120px;
  font-size: 16px;
  font-weight: 700;
  color: #000;
  transition: color 0.3s ease;
}

.user-home.dark-mode .label {
  color: #fff;
}

.value {
  flex: 1;
  font-size: 16px;
  color: #000;
  transition: color 0.3s ease;
}

.user-home.dark-mode .value {
  color: #fff;
}

.not-found {
  text-align: center;
  padding: 60px 20px;
}

.not-found h1 {
  font-size: 32px;
  font-weight: 900;
  color: #000;
  margin-bottom: 20px;
  transition: color 0.3s ease;
}

.user-home.dark-mode .not-found h1 {
  color: #fff;
}

.not-found p {
  font-size: 18px;
  color: #888;
  transition: color 0.3s ease;
}

.user-home.dark-mode .not-found p {
  color: #999;
}

/* 博客列表样式 */
.user-blogs {
  margin-top: 40px;
  padding-top: 30px;
  border-top: 1px solid #000;
  transition: all 0.3s ease;
}

.user-home.dark-mode .user-blogs {
  border-top-color: #fff;
}

.blogs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.blogs-header h2 {
  font-size: 24px;
  font-weight: 900;
  color: #000;
  transition: color 0.3s ease;
}

.user-home.dark-mode .blogs-header h2 {
  color: #fff;
}

.search-box {
  display: flex;
  gap: 10px;
}

.search-box input {
  padding: 8px 12px;
  border: 1px solid #000;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.3s ease;
  background-color: #fff;
  color: #000;
}

.user-home.dark-mode .search-box input {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.search-box button {
  padding: 8px 16px;
  border: 1px solid #000;
  background-color: #000;
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.user-home.dark-mode .search-box button {
  border-color: #fff;
  background-color: #fff;
  color: #000;
}

.search-box button:hover {
  opacity: 0.8;
}

.blogs-content {
  display: flex;
  gap: 40px;
}

/* 左侧时间列表 */
.time-list {
  flex: 0 0 200px;
  background-color: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.user-home.dark-mode .time-list {
  background-color: #1e1e1e;
}

.time-list h3 {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 20px;
  color: #000;
  transition: color 0.3s ease;
}

.user-home.dark-mode .time-list h3 {
  color: #fff;
}

.time-filter {
  margin-bottom: 16px;
}

.time-item {
  padding: 10px;
  cursor: pointer;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
  transition: all 0.3s ease;
  color: #000;
}

.user-home.dark-mode .time-item {
  color: #fff;
}

.time-item:hover {
  background-color: #e0e0e0;
}

.user-home.dark-mode .time-item:hover {
  background-color: #333;
}

.time-item.active {
  background-color: #000;
  color: #fff;
}

.user-home.dark-mode .time-item.active {
  background-color: #fff;
  color: #000;
}

.year-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.year-item {
  border-bottom: 1px solid #ddd;
  padding-bottom: 8px;
  transition: border-color 0.3s ease;
}

.user-home.dark-mode .year-item {
  border-bottom-color: #333;
}

.year-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.year-header:hover {
  background-color: #e0e0e0;
}

.user-home.dark-mode .year-header:hover {
  background-color: #333;
}

.year-header span {
  font-size: 14px;
  font-weight: 500;
  color: #000;
  transition: color 0.3s ease;
}

.user-home.dark-mode .year-header span {
  color: #fff;
}

.expand-icon {
  font-size: 12px;
  transition: transform 0.3s ease;
}

.date-list {
  margin-left: 16px;
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.date-item {
  font-size: 12px;
  padding: 6px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s ease;
  color: #000;
}

.user-home.dark-mode .date-item {
  color: #fff;
}

.date-item:hover {
  background-color: #e0e0e0;
}

.user-home.dark-mode .date-item:hover {
  background-color: #333;
}

.date-item.active {
  background-color: #000;
  color: #fff;
}

.user-home.dark-mode .date-item.active {
  background-color: #fff;
  color: #000;
}

/* 右侧主要内容 */
.main-content {
  flex: 1;
}

.blogs-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.blog-item {
  padding: 20px;
  border: 1px solid #000;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.user-home.dark-mode .blog-item {
  border-color: #fff;
}

.blog-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.05);
}

.user-home.dark-mode .blog-item:hover {
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3);
}

.blog-title {
  font-size: 20px;
  font-weight: 900;
  color: #000;
  margin-bottom: 10px;
  transition: color 0.3s ease;
}

.user-home.dark-mode .blog-title {
  color: #fff;
}

.blog-meta {
  display: flex;
  gap: 20px;
  font-size: 14px;
  color: #888;
  margin-bottom: 10px;
  transition: color 0.3s ease;
}

.user-home.dark-mode .blog-meta {
  color: #999;
}

.blog-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.tag {
  padding: 4px 12px;
  background-color: #f5f5f5;
  border: 1px solid #000;
  border-radius: 16px;
  font-size: 12px;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
}

.user-home.dark-mode .tag {
  background-color: #1e1e1e;
  border-color: #fff;
  color: #fff;
}

.tag:hover {
  background-color: #333;
  color: #fff;
}

.user-home.dark-mode .tag:hover {
  background-color: #ccc;
  color: #000;
}

.blog-preview {
  font-size: 16px;
  line-height: 1.6;
  color: #000;
  margin-bottom: 15px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.3s ease;
}

.user-home.dark-mode .blog-preview {
  color: #fff;
}

.blog-read-more {
  font-size: 14px;
  font-weight: 700;
  color: #000;
  cursor: pointer;
  display: inline-block;
  position: relative;
  transition: color 0.3s ease;
}

.user-home.dark-mode .blog-read-more {
  color: #fff;
}

.blog-read-more::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 100%;
  height: 1px;
  background-color: #000;
  transform: scaleX(0);
  transform-origin: right;
  transition: transform 0.3s ease, background-color 0.3s ease;
}

.user-home.dark-mode .blog-read-more::after {
  background-color: #fff;
}

.blog-read-more:hover::after {
  transform: scaleX(1);
  transform-origin: left;
}

.empty-message {
  text-align: center;
  padding: 40px 20px;
  color: #888;
  transition: color 0.3s ease;
}

.user-home.dark-mode .empty-message {
  color: #999;
}

@media (max-width: 768px) {
  .blogs-content {
    flex-direction: column;
  }
  
  .time-list {
    flex: none;
    width: 100%;
  }
  
  .blogs-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .search-box {
    width: 100%;
  }
  
  .search-box input {
    flex: 1;
  }
  
  .user-home {
    padding: 20px 16px;
  }
  
  .user-profile h1 {
    font-size: 24px;
  }
  
  .user-info h2 {
    font-size: 20px;
  }
  
  .user-blogs h2 {
    font-size: 20px;
  }
  
  .blog-title {
    font-size: 18px;
  }
}

@media (min-width: 768px) {
  .user-home {
    padding: 60px 40px;
  }
  
  .user-profile h1 {
    font-size: 36px;
  }
  
  .user-info h2 {
    font-size: 28px;
  }
  
  .user-blogs h2 {
    font-size: 28px;
  }
  
  .blog-title {
    font-size: 24px;
  }
}

</style>