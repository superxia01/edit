import { Link } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Users, BarChart3, ChevronRight, Shield } from 'lucide-react'

export function AdminPage() {
  const items = [
    {
      title: '用户列表',
      description: '查看所有用户、采集统计与 API Key 状态',
      path: '/admin/users',
      icon: Users,
    },
    {
      title: '全局统计',
      description: '总用户数、总采集量、单篇/批量笔记分布',
      path: '/admin/stats',
      icon: BarChart3,
    },
  ]

  return (
    <main className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2 flex items-center gap-2">
          <Shield className="h-8 w-8" />
          管理后台
        </h1>
        <p className="text-muted-foreground">用户管理、全局统计与系统配置</p>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {items.map((item) => (
          <Link key={item.path} to={item.path}>
            <Card className="h-full transition-colors hover:bg-muted/50 cursor-pointer">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-lg flex items-center gap-2">
                  <item.icon className="h-5 w-5" />
                  {item.title}
                </CardTitle>
                <ChevronRight className="h-5 w-5 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">{item.description}</p>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>
    </main>
  )
}
