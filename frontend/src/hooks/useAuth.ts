import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import apiClient, { adminApi } from '@/api'

interface User {
  id: string
  authCenterUserId: string
  role: string
  nickname?: string
  avatarUrl?: string
  profile?: Record<string, any>
  createdAt: string
  updatedAt: string
}

export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [isAdmin, setIsAdmin] = useState(false)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = async () => {
    const token = localStorage.getItem('token')

    if (!token) {
      setLoading(false)
      return
    }

    try {
      const [meRes, adminRes] = await Promise.all([
        apiClient.get('/user/me'),
        adminApi.checkAdmin().catch(() => ({ isAdmin: false })),
      ])
      const userData = meRes.data?.data || meRes.data
      if (userData) {
        setUser(userData)
      }
      if (adminRes?.isAdmin) {
        setIsAdmin(true)
      }
    } catch (error) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    } finally {
      setLoading(false)
    }
  }

  const login = async (token: string) => {
    localStorage.setItem('token', token)
    await checkAuth()
  }

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setUser(null)
    navigate('/login')
  }

  return {
    user,
    isAdmin,
    loading,
    isAuthenticated: !!user,
    login,
    logout,
    checkAuth,
  }
}
