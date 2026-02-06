import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'

export function AuthCallbackPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading')
  const [message, setMessage] = useState('正在处理登录...')

  useEffect(() => {
    const handleCallback = async () => {
      const token = searchParams.get('token')

      // 统一 token 模式：所有登录方式都返回 token
      if (!token) {
        setStatus('error')
        setMessage('登录失败：缺少 token')
        setTimeout(() => navigate('/login'), 2000)
        return
      }

      try {
        // 用 token 获取用户信息
        const response = await fetch('/api/v1/user/me', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })

        const data = await response.json()

        if (!response.ok || !data.success) {
          setStatus('error')
          setMessage('获取用户信息失败')
          setTimeout(() => navigate('/login'), 2000)
          return
        }

        // 保存 token 和用户信息
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(data.data))

        setStatus('success')
        setMessage('登录成功，正在跳转...')

        // 跳转到 dashboard
        setTimeout(() => navigate('/dashboard'), 1000)
      } catch (error) {
        console.error('登录失败:', error)
        setStatus('error')
        setMessage('登录失败，请稍后重试')
        setTimeout(() => navigate('/login'), 2000)
      }
    }

    handleCallback()
  }, [searchParams, navigate])

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="text-center">
        {status === 'loading' && (
          <div className="space-y-4">
            <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto" />
            <p className="text-lg">{message}</p>
          </div>
        )}

        {status === 'success' && (
          <div className="space-y-4">
            <div className="w-16 h-16 bg-green-500 rounded-full flex items-center justify-center mx-auto">
              <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <p className="text-lg text-green-600">{message}</p>
          </div>
        )}

        {status === 'error' && (
          <div className="space-y-4">
            <div className="w-16 h-16 bg-destructive rounded-full flex items-center justify-center mx-auto">
              <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
            <p className="text-lg text-destructive">{message}</p>
            <p className="text-muted-foreground text-sm mt-2">正在返回登录页...</p>
          </div>
        )}
      </div>
    </div>
  )
}
