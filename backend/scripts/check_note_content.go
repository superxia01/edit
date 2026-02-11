package main

import (
	"fmt"
	"log"

	"github.com/keenchase/edit-business/internal/config"
	"github.com/keenchase/edit-business/internal/model"
	"github.com/keenchase/edit-business/pkg/database"
)

func main() {
	cfg := config.LoadConfig()
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer database.CloseDatabase()

	db := database.GetDB()

	var total, singleCount, batchCount, withContent, singleWithContent int64
	db.Model(&model.Note{}).Count(&total)
	db.Model(&model.Note{}).Where("source = ?", "single").Count(&singleCount)
	db.Model(&model.Note{}).Where("source = ?", "batch").Count(&batchCount)
	db.Model(&model.Note{}).Where("content IS NOT NULL AND TRIM(content) != ''").Count(&withContent)
	db.Model(&model.Note{}).Where("source = ? AND content IS NOT NULL AND TRIM(content) != ''", "single").Count(&singleWithContent)

	fmt.Println("=== notes 表 content 统计 ===")
	fmt.Printf("总笔记数: %d\n", total)
	fmt.Printf("单篇采集(source=single): %d\n", singleCount)
	fmt.Printf("批量采集(source=batch): %d\n", batchCount)
	fmt.Printf("有正文的笔记: %d\n", withContent)
	fmt.Printf("单篇采集且有正文: %d\n", singleWithContent)
	fmt.Println()

	// 取几条单篇且有正文的样例
	var samples []model.Note
	err := db.Where("source = ? AND content IS NOT NULL AND TRIM(content) != ''", "single").
		Order("created_at DESC").
		Limit(3).
		Find(&samples).Error
	if err != nil {
		log.Printf("Query samples error: %v", err)
		return
	}

	if len(samples) == 0 {
		fmt.Println("没有找到单篇采集且有正文的笔记。")
		return
	}

	fmt.Println("=== 单篇且有正文的笔记样例（最近3条）===")
	for i, n := range samples {
		preview := n.Content
		if len(preview) > 80 {
			preview = preview[:80] + "..."
		}
		fmt.Printf("[%d] id=%s title=%s content_len=%d content_preview=%q\n", i+1, n.ID, n.Title, len(n.Content), preview)
	}
}
