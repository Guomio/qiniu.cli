```bash
git clone https://github.com/guomio/qiniu.cli.git
```
```go
// qiniul.cli/qiniu.go
// 修改以下参数为自己七牛云空间的配置

const (
	accessKey = "ecXZ-tsRllsEO6LRu4-Hd9sxxxxxxxx4ZqHPKkwm"
	secretKey = "pktAJPtQp_j4EpozzWx9XPzxxxxxxxxVbKTiAUsa"
	bucket    = "hexo"
	origin    = "http://xxxxxxxxx.bkt.clouddn.com/"
	expires   = 7200
)
```

```bash
go build qiniu.go
ln -s ./qiniu /usr/local/bin/

// 查看使用方法
qiniu --help
```
