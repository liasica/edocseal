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

type documentTask struct {
}

func NewDocumentTask() *documentTask {
	return &documentTask{}
}

func (t *documentTask) Start() {
	c := cron.New()
	_, err := c.AddFunc("00 02 * * * ", func() {
		zap.L().Info("开始执行Document定时删除任务")
		t.Do()
	})
	if err != nil {
		zap.L().Fatal("Document定时删除任务执行失败", zap.Error(err))
		return
	}
	c.Start()
}

func (t *documentTask) Do() {
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

	// 查找已经属于过期时间的合同文档数据
	docs, err := ent.NewDatabase().Document.
		Query().Where(documenEnt.ExpiresAtLTE(time.Now())).
		All(context.Background())
	if err != nil {
		zap.L().Error("查找过期文档数据失败失败", zap.Error(err))
		return
	}

	if len(docs) == 0 {
		zap.L().Info("当前没有已过期的合同文档数据")
		return
	}

	ctx := context.Background()
	err = ent.WithTx(ctx, func(tx *ent.Tx) (err error) {
		// 处理得到过期的文档oss链接以及处理的文档文件夹链接
		var docOssUrls, docOssFolds []string
		for _, doc := range docs {
			// 先删除过期文档
			err = tx.Document.DeleteOne(doc).Exec(ctx)
			if err != nil {
				return
			}

			// 轮询过期的文档数据得到其未签约ossurl
			docPath := doc.Paths
			docOssUrls = append(docOssUrls, docPath.OssUnSigned)
			// 判断是否已签约决定是否删除文件夹
			if doc.Status == documenEnt.StatusUnsigned && doc.SignedURL == "" {
				// 未签约过的数据连同文件夹直接删,分隔未签约url得到其文件夹路径
				usPath := docPath.OssUnSigned
				ossFolder := usPath[0:strings.LastIndex(usPath, "/")]
				docOssFolds = append(docOssFolds, ossFolder)
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
