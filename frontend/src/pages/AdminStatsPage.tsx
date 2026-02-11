import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { adminApi, type AdminStatsOverview } from '@/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { BarChart3, Users, FileText, UserCheck, ArrowLeft } from 'lucide-react'

export function AdminStatsPage() {
  const [overview, setOverview] = useState<AdminStatsOverview | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    loadOverview()
  }, [])

  const loadOverview = async () => {
    try {
      setLoading(true)
      const res = await adminApi.getStatsOverview()
      if (res.data) {
        setOverview(res.data)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载统计失败')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <main className="container mx-auto px-4 py-8">
        <div className="text-center text-muted-foreground">加载中...</div>
      </main>
    )
  }

  return (
    <main className="container mx-auto px-4 py-8">
      <div className="mb-6">
        <Button variant="ghost" size="sm" asChild>
          <Link to="/admin" className="flex items-center gap-1">
            <ArrowLeft className="h-4 w-4" />
            返回管理后台
          </Link>
        </Button>
      </div>

      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2 flex items-center gap-2">
          <BarChart3 className="h-8 w-8" />
          全局统计
        </h1>
        <p className="text-muted-foreground">总用户数、总采集量等</p>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-destructive/10 text-destructive rounded-lg">{error}</div>
      )}

      <div className="grid gap-6 md:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Users className="h-5 w-5" />
              总用户数
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{overview?.totalUsers ?? 0}</p>
            <p className="text-sm text-muted-foreground">已注册用户</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <FileText className="h-5 w-5" />
              总笔记数
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{overview?.totalNotes ?? 0}</p>
            <p className="text-sm text-muted-foreground">已采集笔记</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <UserCheck className="h-5 w-5" />
              总博主数
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{overview?.totalBloggers ?? 0}</p>
            <p className="text-sm text-muted-foreground">已采集博主</p>
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
