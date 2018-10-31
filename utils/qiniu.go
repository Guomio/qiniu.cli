package utils

import (
	"context"
	"path/filepath"
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

// QnConfig 七牛实例配置
type QnConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Origin    string
	Expires   uint32
}

// Qn 七牛实例
type Qn struct {
	QnConfig
	upToken       string
	formUploader  *storage.FormUploader
	bucketManager *storage.BucketManager
	putPolicy     storage.PutPolicy
	born          time.Time
}

// NewQn 创建七牛实例
func NewQn(config *QnConfig) *Qn {
	q := &Qn{}

	q.AccessKey = config.AccessKey
	q.SecretKey = config.SecretKey
	q.Bucket = config.Bucket
	q.Origin = config.Origin
	q.Expires = config.Expires

	q.putPolicy = storage.PutPolicy{
		Scope: config.Bucket,
	}
	q.putPolicy.Expires = config.Expires
	q.born = time.Now()
	mac := qbox.NewMac(config.AccessKey, config.SecretKey)
	q.upToken = q.putPolicy.UploadToken(mac)
	cfgu := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	q.formUploader = storage.NewFormUploader(&cfgu)
	cfgm := storage.Config{
		UseHTTPS: false,
	}
	q.bucketManager = storage.NewBucketManager(mac, &cfgm)
	return q
}

func (q *Qn) freshUpToken() {
	if int(time.Now().Sub(q.born).Seconds()) > int(q.Expires) {
		q.born = time.Now()
		mac := qbox.NewMac(q.AccessKey, q.SecretKey)
		q.upToken = q.putPolicy.UploadToken(mac)
	}
}

// List 获取指定前缀资源列表
func (q *Qn) List(prefix string) ([]storage.ListItem, error) {
	q.freshUpToken()
	limit := 1000
	delimiter := ""
	marker := ""
	task := []storage.ListItem{}
	for {
		entries, _, nextMarker, hashNext, err := q.bucketManager.ListFiles(q.Bucket, prefix, delimiter, marker, limit)
		if err != nil {
			return task, err
		}
		for _, entry := range entries {
			task = append(task, entry)
		}
		if hashNext {
			marker = nextMarker
		} else {
			break
		}
	}
	return task, nil
}

// Upload 上传资源
func (q *Qn) Upload(prefix, localFile, name string) (string, error) {
	q.freshUpToken()
	key := filepath.Join(prefix, time.Now().Format("060102"), name)
	ret := storage.PutRet{}
	err := q.formUploader.PutFile(context.Background(), &ret, q.upToken, key, localFile, nil)
	if err != nil {
		return "上传失败", err
	}
	return ret.Key, nil
}

// Delete 删除资源
func (q *Qn) Delete(keys []string) error {
	q.freshUpToken()
	deleteOps := make([]string, 0, len(keys))
	for _, key := range keys {
		deleteOps = append(deleteOps, storage.URIDelete(q.Bucket, key))
	}
	_, err := q.bucketManager.Batch(deleteOps)
	return err
}

// Rename 重命名资源
func (q *Qn) Rename(srcKey, destKey string) error {
	q.freshUpToken()
	return q.bucketManager.Move(q.Bucket, srcKey, q.Bucket, destKey, true)
}

// Fetch 指定名称抓取网络资源
func (q *Qn) Fetch(restURL, key string) (fileURL string, fetchRetString string, err error) {
	q.freshUpToken()
	fetchRet, err := q.bucketManager.Fetch(restURL, q.Bucket, key)
	fileURL = q.Origin + fetchRet.Key
	fetchRetString = fetchRet.String()
	return
}

// FetchWithoutKey 默认名抓取网络资源
func (q *Qn) FetchWithoutKey(restURL string) (fileURL string, fetchRetString string, err error) {
	q.freshUpToken()
	fetchRet, err := q.bucketManager.FetchWithoutKey(restURL, q.Bucket)
	fileURL = q.Origin + fetchRet.Key
	fetchRetString = fetchRet.String()
	return
}
