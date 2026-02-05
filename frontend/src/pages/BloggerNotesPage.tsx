import { useEffect, useState } from 'react'
import type { ColumnDef } from '@tanstack/react-table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/ui/data-table'
import { notesApi, type Note } from '@/api'
import { ExternalLink, Image as ImageIcon, Video, Trash2 } from 'lucide-react'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'

export function BloggerNotesPage() {
  const [notes, setNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(true)
  const [pagination, setPagination] = useState({
    page: 1,
    size: 20,
    total: 0,
    totalPages: 0,
  })

  const fetchNotes = async (page = 1, size = 20) => {
    try {
      setLoading(true)
      const response = await notesApi.list({
        page,
        size,
        source: 'batch', // 只获取批量采集的笔记
      })
      if (response.data) {
        setNotes(response.data.notes)
        setPagination({
          page: response.data.page,
          size: response.data.size,
          total: response.data.total,
          totalPages: response.data.totalPages,
        })
      }
    } catch (error) {
      console.error('Failed to fetch blogger notes:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchNotes()
  }, [])

  const handleDelete = async (id: string) => {
    try {
      await notesApi.delete(id)
      fetchNotes(pagination.page, pagination.size)
    } catch (error) {
      console.error('Failed to delete note:', error)
    }
  }

  const formatDate = (timestamp: number) => {
    if (!timestamp) return '-'
    return new Date(timestamp).toLocaleDateString('zh-CN')
  }

  const formatNumber = (num: number) => {
    if (num >= 10000) {
      return `${(num / 10000).toFixed(1)}万`
    }
    return num.toString()
  }

  const columns: ColumnDef<Note>[] = [
    {
      accessorKey: 'title',
      header: '标题',
      cell: ({ row }) => (
        <div className="max-w-[300px] truncate" title={row.original.title}>
          {row.original.title || '无标题'}
        </div>
      ),
    },
    {
      accessorKey: 'author',
      header: '作者',
      cell: ({ row }) => row.original.author || '-',
    },
    {
      accessorKey: 'noteType',
      header: '类型',
      cell: ({ row }) => {
        const type = row.original.noteType
        return (
          <Badge variant={type === '视频' ? 'default' : 'secondary'}>
            {type === '视频' ? (
              <Video className="w-3 h-3 mr-1" />
            ) : (
              <ImageIcon className="w-3 h-3 mr-1" />
            )}
            {type || '图文'}
          </Badge>
        )
      },
    },
    {
      accessorKey: 'likes',
      header: '点赞',
      cell: ({ row }) => formatNumber(row.original.likes),
    },
    {
      accessorKey: 'publishDate',
      header: '发布时间',
      cell: ({ row }) => formatDate(row.original.publishDate),
    },
    {
      accessorKey: 'captureTimestamp',
      header: '采集时间',
      cell: ({ row }) => formatDate(row.original.captureTimestamp),
    },
    {
      id: 'actions',
      header: '操作',
      cell: ({ row }) => (
        <div className="flex gap-2">
          <Button variant="ghost" size="sm" asChild>
            <a
              href={row.original.url}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-1"
            >
              查看原文
              <ExternalLink className="w-3 h-3" />
            </a>
          </Button>
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="ghost" size="sm" className="text-destructive">
                <Trash2 className="w-4 h-4" />
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>确认删除</AlertDialogTitle>
                <AlertDialogDescription>
                  确定要删除这条笔记吗？此操作无法撤销。
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>取消</AlertDialogCancel>
                <AlertDialogAction
                  onClick={() => handleDelete(row.original.id)}
                  className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                >
                  删除
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      ),
    },
  ]

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-muted-foreground">加载中...</div>
      </div>
    )
  }

  return (
    <div className="container mx-auto py-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">博主笔记</h1>
        <p className="text-muted-foreground mt-1">
          共 {pagination.total} 条批量采集的笔记
        </p>
      </div>

      <DataTable columns={columns} data={notes} pageSize={pagination.size} />
    </div>
  )
}
