import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { adminApi, type AdminUserDetail } from '@/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import {
  ArrowLeft,
  Key,
  BarChart3,
  Settings,
  Plus,
  Loader2,
  Check,
} from 'lucide-react'

type ApiKeyItem = { id: string; name: string; isActive: boolean; lastUsed?: string; expiresAt?: string | null; createdAt: string }

function APIKeyRow({ apiKey, userId, onUpdated }: { apiKey: ApiKeyItem; userId: string; onUpdated: () => void }) {
  const [expiresInDays, setExpiresInDays] = useState('')
  const [updating, setUpdating] = useState(false)

  const handleUpdateExpiry = async (days: number | null) => {
    try {
      setUpdating(true)
      await adminApi.updateApiKeyExpiry(userId, apiKey.id, { expiresIn: days ?? undefined })
      onUpdated()
    } catch {
      // 错误由父组件展示
    } finally {
      setUpdating(false)
    }
  }

  return (
    <div className="p-4 bg-muted rounded-lg space-y-3">
      <div className="flex items-center justify-between">
        <div>
          <p className="font-medium">{apiKey.name}</p>
          <p className="text-xs text-muted-foreground">
            创建于 {apiKey.createdAt}
            {apiKey.lastUsed && ` · 最后使用 ${apiKey.lastUsed}`}
          </p>
          <p className="text-sm mt-1">
            到期日：<span className="font-medium">{apiKey.expiresAt ?? '永不过期'}</span>
          </p>
        </div>
        <Badge>已激活</Badge>
      </div>
      <div className="flex items-center gap-2 pt-2 border-t">
        <Input
          type="number"
          placeholder="天数"
          className="w-24"
          min={1}
          value={expiresInDays}
          onChange={(e) => setExpiresInDays(e.target.value)}
        />
        <Button
          size="sm"
          variant="outline"
          onClick={() => {
            const d = parseInt(expiresInDays, 10)
            if (!isNaN(d) && d > 0) handleUpdateExpiry(d)
          }}
          disabled={updating}
        >
          {updating ? <Loader2 className="h-4 w-4 animate-spin" /> : '更新有效期'}
        </Button>
        <Button
          size="sm"
          variant="ghost"
          onClick={() => handleUpdateExpiry(null)}
          disabled={updating}
        >
          设为永不过期
        </Button>
      </div>
    </div>
  )
}

export function AdminUserDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [detail, setDetail] = useState<AdminUserDetail | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [creatingKey, setCreatingKey] = useState(false)
  const [newKey, setNewKey] = useState<string | null>(null)
  const [dailyLimit, setDailyLimit] = useState('')
  const [batchLimit, setBatchLimit] = useState('')
  const [apiKeyExpiresIn, setApiKeyExpiresIn] = useState<string>('') // 空=永不过期，数字=天数
  const [saving, setSaving] = useState(false)
  const [saveOk, setSaveOk] = useState(false)

  useEffect(() => {
    if (id) loadDetail()
  }, [id])

  useEffect(() => {
    if (detail?.settings) {
      setDailyLimit(String(detail.settings.collectionDailyLimit))
      setBatchLimit(String(detail.settings.collectionBatchLimit))
    }
  }, [detail?.settings])

  const loadDetail = async () => {
    if (!id) return
    try {
      setLoading(true)
      const res = await adminApi.getUserDetail(id)
      if (res.data) {
        setDetail(res.data)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载用户详情失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateApiKey = async () => {
    if (!id) return
    try {
      setCreatingKey(true)
      setNewKey(null)
      const expiresIn = apiKeyExpiresIn ? parseInt(apiKeyExpiresIn, 10) : undefined
      const res = await adminApi.createApiKeyForUser(id, expiresIn && expiresIn > 0 ? { expiresIn } : undefined)
      if (res.data?.key) {
        setNewKey(res.data.key)
        await loadDetail()
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '创建 API Key 失败')
    } finally {
      setCreatingKey(false)
    }
  }

  const handleSaveSettings = async () => {
    if (!id) return
    const dl = parseInt(dailyLimit, 10)
    const bl = parseInt(batchLimit, 10)
    if (isNaN(dl) || isNaN(bl) || dl < 0 || bl < 0) {
      setError('请输入有效的数字')
      return
    }
    try {
      setSaving(true)
      setError('')
      await adminApi.updateUserSettings(id, {
        collectionDailyLimit: dl,
        collectionBatchLimit: bl,
      })
      setSaveOk(true)
      setTimeout(() => setSaveOk(false), 2000)
      await loadDetail()
    } catch (err) {
      setError(err instanceof Error ? err.message : '更新失败')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <main className="container mx-auto px-4 py-8">
        <div className="text-center text-muted-foreground">加载中...</div>
      </main>
    )
  }

  if (!detail) {
    return (
      <main className="container mx-auto px-4 py-8">
        <div className="text-destructive">{error || '用户不存在'}</div>
        <Button variant="link" asChild>
          <Link to="/admin/users">返回用户列表</Link>
        </Button>
      </main>
    )
  }

  const u = detail.user
  const stats = detail.stats
  const apiKeys = detail.apiKeys || []
  const hasActiveKey = apiKeys.some((k) => k.isActive)

  return (
    <main className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-6">
        <Button variant="ghost" size="sm" asChild>
          <Link to="/admin/users" className="flex items-center gap-1">
            <ArrowLeft className="h-4 w-4" />
            返回用户列表
          </Link>
        </Button>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-destructive/10 text-destructive rounded-lg">{error}</div>
      )}

      {/* 用户信息 */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            {u.avatarUrl ? (
              <img
                src={u.avatarUrl}
                alt=""
                className="w-10 h-10 rounded-full object-cover"
              />
            ) : (
              <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center">
                <span className="text-lg font-bold text-muted-foreground">
                  {(u.nickname || u.authCenterUserId)[0]?.toUpperCase()}
                </span>
              </div>
            )}
            {u.nickname || u.authCenterUserId}
          </CardTitle>
        </CardHeader>
        <CardContent className="text-sm text-muted-foreground space-y-1">
          <p>ID: {u.id}</p>
          <p>AuthCenter User ID: {u.authCenterUserId}</p>
          <p>角色: {u.role}</p>
        </CardContent>
      </Card>

      {/* 统计 */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BarChart3 className="h-5 w-5" />
            采集统计
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-muted-foreground text-sm">笔记数</p>
              <p className="text-2xl font-semibold">{stats?.totalNotes ?? 0}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-sm">博主数</p>
              <p className="text-2xl font-semibold">{stats?.totalBloggers ?? 0}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* API Key */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <span className="flex items-center gap-2">
              <Key className="h-5 w-5" />
              API Key
            </span>
            {!hasActiveKey && (
              <div className="flex items-center gap-2">
                <Input
                  type="number"
                  placeholder="有效期（天），空=永不过期"
                  className="w-36"
                  min={1}
                  value={apiKeyExpiresIn}
                  onChange={(e) => setApiKeyExpiresIn(e.target.value)}
                />
                <Button
                  size="sm"
                  onClick={handleCreateApiKey}
                  disabled={creatingKey}
                >
                  {creatingKey ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4" />}
                  {creatingKey ? '创建中...' : '生成 API Key'}
                </Button>
              </div>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {hasActiveKey ? (
            <div className="space-y-3">
              {apiKeys.filter((k) => k.isActive).map((k) => (
                <APIKeyRow
                  key={k.id}
                  apiKey={k}
                  userId={id!}
                  onUpdated={loadDetail}
                />
              ))}
            </div>
          ) : newKey ? (
            <div className="p-4 bg-green-500/10 text-green-700 dark:text-green-400 rounded-lg">
              <p className="font-semibold mb-2">API Key 已生成，请妥善保存（仅显示一次）：</p>
              <code className="block break-all text-sm">{newKey}</code>
            </div>
          ) : (
            <p className="text-muted-foreground">该用户尚未分配 API Key，点击上方按钮生成。</p>
          )}
        </CardContent>
      </Card>

      {/* 采集设置 */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Settings className="h-5 w-5" />
            采集设置
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground mb-4">
            允许数据收藏开关由用户本人在设置页自主控制（用于防止 API Key 泄露后采集垃圾数据）
          </p>
          <div className="space-y-4">
            <div className="flex items-center gap-4">
              <div className="flex-1">
                <Label htmlFor="dailyLimit">每日限额</Label>
                <Input
                  id="dailyLimit"
                  type="number"
                  min={0}
                  value={dailyLimit}
                  onChange={(e) => setDailyLimit(e.target.value)}
                />
              </div>
              <div className="flex-1">
                <Label htmlFor="batchLimit">单次限额</Label>
                <Input
                  id="batchLimit"
                  type="number"
                  min={0}
                  value={batchLimit}
                  onChange={(e) => setBatchLimit(e.target.value)}
                />
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Button onClick={handleSaveSettings} disabled={saving}>
                {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Check className="h-4 w-4" />}
                {saving ? '保存中...' : '保存'}
              </Button>
              {saveOk && (
                <span className="text-sm text-green-600">已保存</span>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </main>
  )
}
