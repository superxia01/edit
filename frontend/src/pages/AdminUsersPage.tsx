import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { adminApi, type AdminUserListItem } from '@/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Users, Key, ChevronRight, ChevronLeft } from 'lucide-react'

const PAGE_SIZE = 20

export function AdminUsersPage() {
  const [users, setUsers] = useState<AdminUserListItem[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    loadUsers(page)
  }, [page])

  const loadUsers = async (p: number) => {
    try {
      setLoading(true)
      const res = await adminApi.listUsers({ page: p, size: PAGE_SIZE })
      if (res.data) {
        setUsers(res.data.items)
        setTotal(res.data.total)
        setTotalPages(res.data.totalPages)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载用户列表失败')
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
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold mb-2 flex items-center gap-2">
            <Users className="h-8 w-8" />
            用户管理
          </h1>
          <p className="text-muted-foreground">查看所有用户、采集统计与 API Key 状态</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" asChild>
            <Link to="/admin">返回管理后台</Link>
          </Button>
          <Button variant="outline" asChild>
            <Link to="/admin/stats">全局统计</Link>
          </Button>
        </div>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-destructive/10 text-destructive rounded-lg">{error}</div>
      )}

      <Card>
        <CardHeader>
          <CardTitle>用户列表</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>用户</TableHead>
                <TableHead>AuthCenter ID</TableHead>
                <TableHead>笔记数</TableHead>
                <TableHead>博主数</TableHead>
                <TableHead>API Key</TableHead>
                <TableHead>注册时间</TableHead>
                <TableHead className="w-[80px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center text-muted-foreground">
                    暂无用户
                  </TableCell>
                </TableRow>
              ) : (
                users.map((u) => (
                  <TableRow key={u.id}>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        {u.avatarUrl ? (
                          <img
                            src={u.avatarUrl}
                            alt=""
                            className="w-8 h-8 rounded-full object-cover"
                          />
                        ) : (
                          <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                            <Users className="h-4 w-4 text-muted-foreground" />
                          </div>
                        )}
                        <span>{u.nickname || u.authCenterUserId.slice(0, 8) + '...'}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <code className="text-xs">{u.authCenterUserId}</code>
                    </TableCell>
                    <TableCell>{u.totalNotes}</TableCell>
                    <TableCell>{u.totalBloggers}</TableCell>
                    <TableCell>
                      {u.hasApiKey ? (
                        <Badge variant="default" className="gap-1">
                          <Key className="h-3 w-3" />
                          已有
                        </Badge>
                      ) : (
                        <Badge variant="secondary">未分配</Badge>
                      )}
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {u.createdAt}
                    </TableCell>
                    <TableCell>
                      <Button variant="ghost" size="sm" asChild>
                        <Link to={`/admin/users/${u.id}`}>
                          <ChevronRight className="h-4 w-4" />
                        </Link>
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
          {totalPages > 1 && (
            <div className="flex items-center justify-between mt-4">
              <p className="text-sm text-muted-foreground">
                共 {total} 条，第 {page}/{totalPages} 页
              </p>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page <= 1}
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                >
                  <ChevronLeft className="h-4 w-4" />
                  上一页
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page >= totalPages}
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                >
                  下一页
                  <ChevronRight className="h-4 w-4" />
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </main>
  )
}
