import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'

export function AuthCallbackPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading')
  const [message, setMessage] = useState('正在处理登录...')

  useEffect(() => {
    const handleCallback = async () => {
      const code = searchParams.get('code')
      const type = searchParams.get('type')
      const token = searchParams.get('token')
      const userId = searchParams.get('userId')

      // 微信内登录：auth-center 直接返回了 token 和 userId
      if (token && userId) {
        try {
          // 用 token 获取用户信息
          const response = await fetch('/api/v1/users/me', {
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
          return
        } catch (error) {
          console.error('获取用户信息失败:', error)
          setStatus('error')
          setMessage('登录失败，请稍后重试')
          setTimeout(() => navigate('/login'), 2000)
          return
        }
      }

      // PC扫码登录：需要用 code 换取 token
      if (!code) {
        setStatus('error')
        setMessage('登录失败：缺少必要参数')
        setTimeout(() => navigate('/login'), 2000)
        return
      }

      try {
        // 调用后端 API，用 code 换取用户信息
        const response = await fetch(`/api/v1/auth/wechat/callback?code=${encodeURIComponent(code)}&type=${encodeURIComponent(type || 'open')}`)

        const data = await response.json()

        // 调试日志
        console.log('[AuthCallback] Response:', data)
        console.log('[AuthCallback] Response OK:', response.ok)
        console.log('[AuthCallback] Has data.data:', !!data.data)

        // 兼容两种响应格式：success 字段或 code === 0
        const isSuccess = data.success || data.code === 0

        if (!response.ok || !isSuccess) {
          console.error('[AuthCallback] Login failed:', data)
          setStatus('error')
          setMessage(data.message || '登录失败')
          setTimeout(() => navigate('/login'), 2000)
          return
        }

        // 保存 token 和用户信息
        console.log('[AuthCallback] Saving token:', data.data?.token?.substring(0, 50))
        localStorage.setItem('token', data.data.token)
        localStorage.setItem('user', JSON.stringify(data.data.user))
        console.log('[AuthCallback] Token saved successfully')

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
