package task

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/liasica/edocseal"
	"github.com/liasica/edocseal/internal/ent"
	documenEnt "github.com/liasica/edocseal/internal/ent/document"
	"github.com/liasica/edocseal/internal/g"
)

type FileTask struct {
}

func NewFileTask() *FileTask {
	return &FileTask{}
}

func (t *FileTask) Run() {
	c := cron.New()
	_, err := c.AddFunc("0 2 * * *", func() {
		zap.L().Info("开始执行Document定时删除任务")
		t.do()
	})
	if err != nil {
		zap.L().Fatal("Document定时删除任务执行失败", zap.Error(err))
		return
	}
	c.Start()
}

func (t *FileTask) do() {
	// 删除本地上月合同文件夹
	lastMonth := time.Now().AddDate(0, -1, 0).Format("200601")
	delMonthDir, _ := filepath.Abs(filepath.Join(g.GetDocumentDir(), lastMonth[:4], lastMonth[4:6]))
	err := edocseal.DeleteDirectory(delMonthDir)
	if err != nil {
		zap.L().Error("定时删除上月合同文件夹失败", zap.Error(err))
		return
	}
	// 删除本地前一天合同数据
	lastDay := time.Now().AddDate(0, 0, -1).Format("20060102")
	delDayDir, _ := filepath.Abs(filepath.Join(g.GetDocumentDir(), lastDay[:4], lastDay[4:6], lastDay[6:]))
	err = edocseal.DeleteDirectory(delDayDir)
	if err != nil {
		zap.L().Error("定时删除昨日合同文件夹失败", zap.Error(err))
		return
	}

	// 查询前一天数据
	yesterdayBeginTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().AddDate(0, 0, -1).Day(), 0, 0, 0, 0, time.Local)
	yesterdayEndTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	var docs []*ent.Document
	docs, err = ent.NewDatabase().Document.Query().
		Where(
			documenEnt.CreateAtGTE(yesterdayBeginTime),
			documenEnt.CreateAtLT(yesterdayEndTime),
		).
		All(context.Background())
	if err != nil {
		zap.L().Error("查找昨日文档数据失败失败", zap.Error(err))
		return
	}

	if len(docs) == 0 {
		zap.L().Info("昨日没有合同文档数据")
		return
	}

	ctx := context.Background()
	err = ent.WithTx(ctx, func(tx *ent.Tx) (err error) {
		// 处理得到过期的文档oss链接以及处理的文档文件夹链接
		var docOssUrls, docOssFolds []string

		for _, doc := range docs {
			docPath := doc.Paths
			// 判断是否已签约
			switch doc.Status {
			case documenEnt.StatusUnsigned:
				// 未签约且已过期数据删除本地数据库记录、oss云端临时未签约文件、oss云端文件所属文件夹
				// 前一天未签约数据此时必过期，不加过期时间的判断
				// 先删除过期文档
				err = tx.Document.DeleteOne(doc).Exec(ctx)
				if err != nil {
					return
				}
				// 未签约文档url
				docOssUrls = append(docOssUrls, docPath.OssUnSigned)
				// 未签约文档文件夹url
				ossFolder := docPath.OssUnSigned[0:strings.LastIndex(docPath.OssUnSigned, "/")]
				docOssFolds = append(docOssFolds, ossFolder)
			case documenEnt.StatusSigned:
				// 已签约数据仅删除oss云端临时未签约文件
				// 未签约文档url
				docOssUrls = append(docOssUrls, docPath.OssUnSigned)
			default:
				continue
			}

		}

		// oss多文件删除
		cfg := g.GetAliyunOss()
		ao, err := edocseal.NewAliyunOss(cfg.AccessKeyId, cfg.AccessKeySecret, cfg.Endpoint, cfg.Bucket)
		if err != nil {
			return
		}
		// 删除未签约文件
		if len(docOssUrls) > 0 {
			err = ao.MultiDelete(docOssUrls)
			if err != nil {
				return
			}
		}
		// 删除只包含未签约文件的文件夹
		if len(docOssFolds) > 0 {
			err = ao.MultiDelete(docOssFolds)
			if err != nil {
				return
			}
		}

		return
	})
	if err != nil {
		zap.L().Error("删除文档失败", zap.Error(err))
		return
	}

}
