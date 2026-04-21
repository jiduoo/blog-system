<template>
  <div class="blogs-container" :class="{ 'dark-mode': themeStore.isDarkMode }">
    <div class="blogs-header">
      <h1>博客</h1>
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
        <!-- 标签过滤器 -->
        <div class="tags-filter">
          <button 
            :class="{ 'active': selectedTag === '' && !selectedDate }"
            @click="clearFilters"
          >
            全部标签
          </button>
          <button 
            v-for="tag in tags" 
            :key="tag.ID"
            :class="{ 'active': selectedTag === tag.Name }"
            @click="filterByTag(tag.Name)"
          >
            {{ tag.Name }}
          </button>
        </div>
        
        <div class="blogs-list">
          <div v-for="blog in filteredBlogs" :key="blog.ID" class="blog-item">
            <div class="blog-title">{{ blog.title }}</div>
            <div class="blog-meta">
              <span>{{ blog.author }}</span>
              <span>{{ formatDate(blog.CreatedAt) }}</span>
              <span>{{ blog.views }} 浏览</span>
            </div>
            <div class="blog-tags">
              <span 
                v-for="tag in blog.tags" 
                :key="tag.ID"
                class="tag"
                @click="filterByTag(tag.Name)"
              >
                {{ tag.Name }}
              </span>
            </div>
            <div class="blog-preview">{{ blog.preview }}</div>
            <div class="blog-read-more" @click="navigateToDetail(blog.ID)">阅读更多</div>
          </div>
        </div>
        
        <div v-if="filteredBlogs.length === 0" class="empty-message">
          {{ selectedDate ? '该日期没有文章' : '暂无文章' }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from '../axios';
import { useThemeStore } from '../store/theme';

interface Tag {
  ID: number;
  Name: string;
}

interface Blog {
  ID: number;
  title: string;
  author: string;
  CreatedAt: string;
  views: number;
  preview: string;
  tags: Tag[];
}

interface DateGroup {
  date: string;
  count: number;
}

interface YearGroup {
  year: string;
  dates: DateGroup[];
}

const router = useRouter();
const themeStore = useThemeStore();
const blogs = ref<Blog[]>([]);
const tags = ref<Tag[]>([]);
const selectedTag = ref('');
const selectedDate = ref('');
const expandedYears = ref<string[]>([]);
const searchKeyword = ref('');

const fetchBlogs = async () => {
  try {
    const response = await axios.get('/blogs');
    blogs.value = response.data;
  } catch (error) {
    console.error('获取博客列表失败:', error);
  }
};

const fetchTags = async () => {
  try {
    const response = await axios.get('/tags');
    tags.value = response.data;
  } catch (error) {
    console.error('获取标签列表失败:', error);
  }
};

const filterByTag = async (tagName: string) => {
  selectedTag.value = tagName;
  selectedDate.value = '';
  searchKeyword.value = '';
  
  if (tagName) {
    try {
      const response = await axios.get(`/tag-blogs/${encodeURIComponent(tagName)}`);
      blogs.value = response.data;
    } catch (error) {
      console.error('按标签获取博客失败:', error);
    }
  } else {
    await fetchBlogs();
  }
};

const searchBlogs = async () => {
  if (searchKeyword.value) {
    try {
      const response = await axios.get(`/blogs/search?keyword=${encodeURIComponent(searchKeyword.value)}`);
      blogs.value = response.data;
    } catch (error) {
      console.error('搜索博客失败:', error);
    }
  } else {
    await fetchBlogs();
  }
  selectedTag.value = '';
  selectedDate.value = '';
};

const filterByDate = (date: string) => {
  selectedDate.value = date;
  selectedTag.value = '';
  searchKeyword.value = '';
};

const clearFilters = async () => {
  selectedTag.value = '';
  selectedDate.value = '';
  searchKeyword.value = '';
  await fetchBlogs();
};

const clearDateFilter = async () => {
  selectedDate.value = '';
  if (searchKeyword.value) {
    await searchBlogs();
  } else if (selectedTag.value) {
    await filterByTag(selectedTag.value);
  } else {
    await fetchBlogs();
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
  let result = blogs.value;
  
  if (selectedDate.value) {
    result = result.filter(blog => {
      const blogDate = new Date(blog.CreatedAt);
      const month = (blogDate.getMonth() + 1).toString().padStart(2, '0');
      const day = blogDate.getDate().toString().padStart(2, '0');
      const dateStr = `${month}-${day}`;
      return dateStr === selectedDate.value;
    });
  }
  
  return result;
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
  fetchBlogs();
  fetchTags();
});
</script>

<style scoped>
.blogs-container {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #fff;
  color: #000;
  transition: all 0.3s ease;
}

.blogs-container.dark-mode {
  background-color: #000;
  color: #fff;
}

.blogs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40px;
}

.blogs-header h1 {
  font-size: 36px;
  font-weight: 900;
  color: #000;
  transition: color 0.3s ease;
}

.blogs-container.dark-mode .blogs-header h1 {
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

.blogs-container.dark-mode .search-box input {
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

.blogs-container.dark-mode .search-box button {
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

.blogs-container.dark-mode .time-list {
  background-color: #1e1e1e;
}

.time-list h3 {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 20px;
  color: #000;
  transition: color 0.3s ease;
}

.blogs-container.dark-mode .time-list h3 {
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

.blogs-container.dark-mode .time-item {
  color: #fff;
}

.time-item:hover {
  background-color: #e0e0e0;
}

.blogs-container.dark-mode .time-item:hover {
  background-color: #333;
}

.time-item.active {
  background-color: #000;
  color: #fff;
}

.blogs-container.dark-mode .time-item.active {
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

.blogs-container.dark-mode .year-item {
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

.blogs-container.dark-mode .year-header:hover {
  background-color: #333;
}

.year-header span {
  font-size: 14px;
  font-weight: 500;
  color: #000;
  transition: color 0.3s ease;
}

.blogs-container.dark-mode .year-header span {
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

.blogs-container.dark-mode .date-item {
  color: #fff;
}

.date-item:hover {
  background-color: #e0e0e0;
}

.blogs-container.dark-mode .date-item:hover {
  background-color: #333;
}

.date-item.active {
  background-color: #000;
  color: #fff;
}

.blogs-container.dark-mode .date-item.active {
  background-color: #fff;
  color: #000;
}

/* 右侧主要内容 */
.main-content {
  flex: 1;
}

.tags-filter {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 32px;
}

.tags-filter button {
  padding: 8px 16px;
  border: 1px solid #000;
  background-color: #fff;
  font-size: 14px;
  font-weight: 500;
  color: #000;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 16px;
}

.blogs-container.dark-mode .tags-filter button {
  border-color: #fff;
  background-color: #000;
  color: #fff;
}

.tags-filter button:hover {
  background-color: #333;
  color: #fff;
}

.blogs-container.dark-mode .tags-filter button:hover {
  background-color: #ccc;
  color: #000;
}

.tags-filter button.active {
  background-color: #000;
  color: #fff;
}

.blogs-container.dark-mode .tags-filter button.active {
  background-color: #fff;
  color: #000;
}

.blogs-list {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.blog-item {
  padding: 20px;
  border-bottom: 1px solid #000;
  transition: all 0.3s ease;
}

.blogs-container.dark-mode .blog-item {
  border-bottom-color: #fff;
}

.blog-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.05);
}

.blogs-container.dark-mode .blog-item:hover {
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3);
}

.blog-title {
  font-size: 24px;
  font-weight: 900;
  color: #000;
  margin-bottom: 12px;
  transition: color 0.3s ease;
}

.blogs-container.dark-mode .blog-title {
  color: #fff;
}

.blog-meta {
  display: flex;
  gap: 20px;
  font-size: 14px;
  color: #888;
  margin-bottom: 12px;
  transition: color 0.3s ease;
}

.blogs-container.dark-mode .blog-meta {
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

.blogs-container.dark-mode .tag {
  background-color: #1e1e1e;
  border-color: #fff;
  color: #fff;
}

.tag:hover {
  background-color: #333;
  color: #fff;
}

.blogs-container.dark-mode .tag:hover {
  background-color: #ccc;
  color: #000;
}

.blog-preview {
  font-size: 16px;
  line-height: 1.6;
  color: #000;
  margin-bottom: 20px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.3s ease;
}

.blogs-container.dark-mode .blog-preview {
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

.blogs-container.dark-mode .blog-read-more {
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

.blogs-container.dark-mode .blog-read-more::after {
  background-color: #fff;
}

.blog-read-more:hover::after {
  transform: scaleX(1);
  transform-origin: left;
}

.empty-message {
  text-align: center;
  padding: 40px;
  font-size: 16px;
  color: #888;
}

.blogs-container.dark-mode .empty-message {
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
  
  .blogs-container {
    padding: 20px 16px;
  }
  
  .blog-title {
    font-size: 20px;
  }
}

@media (min-width: 768px) {
  .blogs-container {
    padding: 60px 40px;
  }
  
  .blog-title {
    font-size: 28px;
  }
}
</style>
