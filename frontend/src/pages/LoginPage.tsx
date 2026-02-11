import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { authApi } from '@/api'

export function LoginPage() {
  const [loginType, setLoginType] = useState<'wechat' | 'password'>('wechat')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { login } = useAuth()

  // 从 URL 读取 error 参数
  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const errorMsg = params.get('error')
    if (errorMsg) {
      setError(decodeURIComponent(errorMsg))
      window.history.replaceState({}, '', window.location.pathname)
    }
  }, [])

  const handleWeChatLogin = () => {
    setError('')
    setLoading(true)
    window.location.href = '/api/v1/auth/wechat/login'
  }

  const handlePasswordLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    const formData = new FormData(e.currentTarget)
    const phoneNumber = formData.get('phoneNumber') as string
    const password = formData.get('password') as string

    try {
      const res: any = await authApi.passwordLogin({ phoneNumber, password })
      const token = res?.data?.token
      if (token) {
        await login(token)
        navigate('/dashboard')
      } else {
        setError('登录失败，请重试')
      }
    } catch (err: any) {
      setError(err.response?.data?.message || err.message || '登录失败，请重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4">
      <div className="max-w-md w-full">
        {/* Logo 和标题 */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">内容管理工具</h1>
          <p className="text-gray-600">笔记与创作者数据管理</p>
        </div>

        {/* 登录卡片 */}
        <Card className="shadow-xl">
          <CardHeader className="space-y-4">
            {/* 登录方式切换 */}
            <div className="flex bg-muted rounded-lg p-1">
              <Button
                onClick={() => setLoginType('wechat')}
                variant={loginType === 'wechat' ? 'default' : 'ghost'}
                className="flex-1"
              >
                微信登录
              </Button>
              <Button
                onClick={() => setLoginType('password')}
                variant={loginType === 'password' ? 'default' : 'ghost'}
                className="flex-1"
              >
                手机号登录
              </Button>
            </div>

            {/* 错误提示 */}
            {error && (
              <Badge variant="destructive" className="w-full py-2 px-3 justify-center">
                {error}
              </Badge>
            )}
          </CardHeader>

          <CardContent>
            {/* 微信登录 */}
            {loginType === 'wechat' && (
              <div className="space-y-4">
                <div className="text-center text-muted-foreground mb-6">
                  <p>点击下方按钮使用微信账号登录</p>
                </div>
                <Button
                  onClick={handleWeChatLogin}
                  disabled={loading}
                  className="w-full bg-green-500 hover:bg-green-600"
                  size="lg"
                >
                  <svg className="w-6 h-6 mr-2" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.01-.27-.027-.407-.03zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.969-.982z" />
                  </svg>
                  {loading ? '登录中...' : '微信登录'}
                </Button>
              </div>
            )}

            {/* 密码登录 */}
            {loginType === 'password' && (
              <form onSubmit={handlePasswordLogin} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="phoneNumber">手机号</Label>
                  <Input
                    type="tel"
                    id="phoneNumber"
                    name="phoneNumber"
                    required
                    placeholder="请输入手机号"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="password">密码</Label>
                  <Input
                    type="password"
                    id="password"
                    name="password"
                    required
                    placeholder="请输入密码"
                  />
                </div>

                <Button
                  type="submit"
                  disabled={loading}
                  className="w-full"
                  size="lg"
                >
                  {loading ? '登录中...' : '登录'}
                </Button>
              </form>
            )}
          </CardContent>
        </Card>

        {/* 底部提示 */}
        <p className="text-center text-muted-foreground text-sm mt-6">
          登录即表示同意《用户协议》和《隐私政策》
        </p>
      </div>
    </div>
  )
}
