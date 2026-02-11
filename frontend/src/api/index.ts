import axios from 'axios'

// 生产环境使用相对路径，开发环境可通过 .env 配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 添加 token 到请求头
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    // 处理 401 未授权错误
    if (error.response?.status === 401) {
      console.error('401 Unauthorized - Token may be invalid or expired')
      // 暂时不删除 token，以便调试
      // localStorage.removeItem('token')
      // localStorage.removeItem('user')
      // window.location.href = '/login'
    }
    const message = error.response?.data?.message || error.message || '请求失败'
    console.error('API Error:', message)
    return Promise.reject(new Error(message))
  }
)

// API 响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

// ========== 笔记相关类型 ==========
export interface Note {
  id: string
  url: string
  title: string
  author: string
  content: string
  tags: string[]
  imageUrls: string[]
  videoUrl?: string
  noteType: string
  coverImageUrl: string
  likes: number
  collects: number
  comments: number
  publishDate: number
  source: string // 'single' | 'batch'
  captureTimestamp: number
  createdAt: string
  updatedAt: string
}

export interface CreateNoteRequest {
  url: string
  title?: string
  author?: string
  content?: string
  tags?: string[]
  imageUrls?: string[]
  videoUrl?: string
  noteType?: string
  coverImageUrl?: string
  likes?: number
  collects?: number
  comments?: number
  publishDate?: number
  captureTimestamp: number
}

export interface ListNotesRequest {
  page?: number
  size?: number
  author?: string
  tags?: string[]
  source?: string // 'single' | 'batch'
}

export interface ListNotesResponse {
  notes: Note[]
  total: number
  page: number
  size: number
  totalPages: number
}

// ========== 博主相关类型 ==========
export interface Blogger {
  id: string
  xhsId: string
  bloggerName: string
  avatarUrl: string
  description: string
  followersCount: number
  bloggerUrl: string
  captureTimestamp: number
  createdAt: string
  updatedAt: string
}

export interface CreateBloggerRequest {
  xhsId: string
  bloggerName?: string
  avatarUrl?: string
  description?: string
  followersCount?: number
  bloggerUrl?: string
  captureTimestamp: number
}

export interface ListBloggersRequest {
  page?: number
  size?: number
}

export interface ListBloggersResponse {
  bloggers: Blogger[]
  total: number
  page: number
  size: number
  totalPages: number
}

// ========== 统计数据类型 ==========
export interface Stats {
  totalNotes: number
  totalBloggers: number
  totalLikes: number
  totalCollects: number
  totalComments: number
  imageNotes: number
  videoNotes: number
}

// ========== 笔记 API ==========
export const notesApi = {
  // 创建笔记
  create: (data: CreateNoteRequest) =>
    apiClient.post<any, ApiResponse<Note>>('/notes', data),

  // 批量创建笔记
  batchCreate: (data: CreateNoteRequest[]) =>
    apiClient.post<any, ApiResponse>('/notes/batch', data),

  // 获取笔记列表
  list: (params: ListNotesRequest) =>
    apiClient.get<any, ApiResponse<ListNotesResponse>>('/notes', { params }),

  // 获取笔记详情
  getById: (id: string) =>
    apiClient.get<any, ApiResponse<Note>>(`/notes/${id}`),

  // 更新笔记
  update: (id: string, data: Partial<Note>) =>
    apiClient.put<any, ApiResponse<Note>>(`/notes/${id}`, data),

  // 删除笔记
  delete: (id: string) =>
    apiClient.delete<any, ApiResponse>(`/notes/${id}`),
}

// ========== 博主 API ==========
export const bloggersApi = {
  // 创建博主
  create: (data: CreateBloggerRequest) =>
    apiClient.post<any, ApiResponse<Blogger>>('/bloggers', data),

  // 批量创建博主
  batchCreate: (data: CreateBloggerRequest[]) =>
    apiClient.post<any, ApiResponse>('/bloggers/batch', data),

  // 插入或更新博主
  upsert: (data: CreateBloggerRequest) =>
    apiClient.post<any, ApiResponse<Blogger>>('/bloggers/upsert', data),

  // 获取博主列表
  list: (params: ListBloggersRequest) =>
    apiClient.get<any, ApiResponse<ListBloggersResponse>>('/bloggers', { params }),

  // 获取博主详情
  getById: (id: string) =>
    apiClient.get<any, ApiResponse<Blogger>>(`/bloggers/${id}`),

  // 根据 xhs_id 获取博主
  getByXhsId: (xhsId: string) =>
    apiClient.get<any, ApiResponse<Blogger>>(`/bloggers/xhs/${xhsId}`),

  // 更新博主
  update: (id: string, data: Partial<Blogger>) =>
    apiClient.put<any, ApiResponse<Blogger>>(`/bloggers/${id}`, data),

  // 删除博主
  delete: (id: string) =>
    apiClient.delete<any, ApiResponse>(`/bloggers/${id}`),
}

// ========== 统计数据 API ==========
export const statsApi = {
  // 获取统计数据
  get: () =>
    apiClient.get<any, ApiResponse<Stats>>('/stats'),
}

// ========== API Key 相关类型 ==========
export interface APIKey {
  id: string
  name: string
  key: string
  isActive: boolean
  lastUsed: string | null
  expiresAt: string | null
  createdAt: string
}

export interface APIKeyStats {
  totalCount: number
  activeCount: number
  totalUsage: number
  lastUsed: string | null
}

export interface CreateAPIKeyRequest {
  name: string
  expiresIn?: number // Optional expiration in days
}

// ========== API Key API ==========
export const apiKeyApi = {
  // 获取 API Key（不再自动创建，无则 404）
  getOrCreate: () =>
    apiClient.get<any, ApiResponse<APIKey>>('/api-keys/get-or-create'),
}

// ========== Admin 相关类型 ==========
export interface AdminUserListItem {
  id: string
  authCenterUserId: string
  nickname?: string
  avatarUrl?: string
  role: string
  createdAt: string
  totalNotes: number
  totalBloggers: number
  hasApiKey: boolean
}

export interface AdminUserDetail {
  user: {
    id: string
    authCenterUserId: string
    nickname?: string
    avatarUrl?: string
    role: string
    createdAt: string
  }
  stats: {
    totalNotes: number
    totalBloggers: number
  }
  settings: {
    userId: string
    collectionEnabled: boolean
    collectionDailyLimit: number
    collectionBatchLimit: number
  }
  apiKeys: Array<{ id: string; name: string; isActive: boolean; lastUsed?: string; expiresAt?: string | null; createdAt: string }>
}

export interface AdminListUsersResponse {
  items: AdminUserListItem[]
  total: number
  page: number
  size: number
  totalPages: number
}

export interface AdminStatsOverview {
  totalUsers: number
  totalNotes: number
  totalBloggers: number
}

// ========== Admin API ==========
export const adminApi = {
  checkAdmin: () =>
    apiClient.get<any, { isAdmin: boolean }>('/admin/check'),

  listUsers: (params?: { page?: number; size?: number }) =>
    apiClient.get<any, ApiResponse<AdminListUsersResponse>>('/admin/users', { params }),

  getUserDetail: (userId: string) =>
    apiClient.get<any, ApiResponse<AdminUserDetail>>(`/admin/users/${userId}`),

  createApiKeyForUser: (userId: string, data?: { expiresIn?: number }) =>
    apiClient.post<any, ApiResponse<APIKey>>(`/admin/users/${userId}/api-keys`, data ?? {}),

  updateApiKeyExpiry: (userId: string, apiKeyId: string, data?: { expiresIn?: number | null }) =>
    apiClient.patch<any, ApiResponse>(`/admin/users/${userId}/api-keys/${apiKeyId}/expiry`, data ?? {}),

  updateUserSettings: (userId: string, data: {
    collectionDailyLimit?: number
    collectionBatchLimit?: number
    collectionEnabled?: boolean
  }) =>
    apiClient.put<any, ApiResponse>(`/admin/users/${userId}/settings`, data),

  getStatsOverview: () =>
    apiClient.get<any, ApiResponse<AdminStatsOverview>>('/admin/stats/overview'),
}

// ========== User Settings 相关类型 ==========
export interface UserSettings {
  userId: string
  collectionEnabled: boolean
  collectionDailyLimit: number
  collectionBatchLimit: number
}

export interface ToggleCollectionRequest {
  enabled: boolean
}

// ========== User Settings API ==========
export const userSettingsApi = {
  // 获取或创建用户设置
  getOrCreate: () =>
    apiClient.get<any, ApiResponse<UserSettings>>('/user-settings'),

  // 切换采集开关
  toggleCollection: (enabled: boolean) =>
    apiClient.post<any, ApiResponse<UserSettings>>('/user-settings/toggle-collection', { enabled }),
}

// ========== 认证 API ==========
export interface PasswordLoginRequest {
  phoneNumber: string
  password: string
}

export interface PasswordLoginResponse {
  success: boolean
  data?: { token: string }
  message?: string
}

export const authApi = {
  passwordLogin: (data: PasswordLoginRequest) =>
    apiClient.post<any, PasswordLoginResponse>('/auth/password', data),
}

export default apiClient
export { apiClient }

