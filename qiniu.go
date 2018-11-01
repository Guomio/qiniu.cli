package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/guomio/qiniu.cli/utils"

	"github.com/urfave/cli"
)

var qn *utils.Qn

const (
	accessKey = "ecXZ-tsRllsEO6LRu4-Hd9sxxxxxxx14ZqHPKkwm"
	secretKey = "pktAJPtQp_j4EpozzWx9XPzxxxxxxx0VbKTiAUsa"
	bucket    = "hexo"
	origin    = "http://xxxxxxxx.bkt.clouddn.com/"
	expires   = 7200
)

func init() {
	qnConfig := &utils.QnConfig{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Origin:    origin,
		Expires:   expires,
	}
	qn = utils.NewQn(qnConfig)
}

func main() {
	app := cli.NewApp()
	app.Name = "qiniu"
	app.Usage = "七牛云空间文件管理"
	app.Commands = appCommands()
	app.Action = appAction

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func appCommands() []cli.Command {
	return []cli.Command{
		cli.Command{
			Name:    "upload",
			Aliases: []string{"u"},
			Usage:   "上传一个单文件",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "prefix, p",
					Value: "hexo",
					Usage: "输入文件上传前缀",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Println("请输入要上传的文件路径")
					return nil
				}
				for i, arg := range c.Args() {
					fileInfo, err := os.Stat(utils.GetAbsPath(arg))
					if err != nil {
						return err
					}
					fmt.Printf("正在上传 [%d/%d] 个文件，大小: %s\n", i+1, len(c.Args()), utils.GetFileSize(fileInfo.Size()))
					if fileInfo.IsDir() {
						fmt.Printf("%s 是文件夹，跳过上传\n", fileInfo.Name())
						continue
					}
					filePath, err := qn.Upload(c.String("prefix"), utils.GetAbsPath(arg), fileInfo.Name())
					if err == nil {
						fmt.Println(qn.Origin + filePath)
					} else {
						fmt.Printf("文件上传失败: %s\n", arg)
					}
				}
				return nil
			},
		},
		cli.Command{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "通过文件名称删除指定文件，支持输入多个，以空格分割",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Println("请输入要删除的文件名")
					return nil
				}
				return qn.Delete(c.Args())
			},
		},
		cli.Command{
			Name:    "rename",
			Aliases: []string{"r"},
			Usage:   "重命名文件名，输入源文件名和目标文件名，以空格分割",
			Action: func(c *cli.Context) error {
				if c.NArg() < 2 {
					fmt.Println("请正确输入文件名和新文件名")
					return nil
				}
				args := c.Args()
				return qn.Rename(args[0], args[1])
			},
		},
		cli.Command{
			Name:    "fetch",
			Aliases: []string{"f"},
			Usage:   "抓取网络资源到空间，一个参数时以默认名存储，两个参数时，以第二个参数为文件名存储",
			Action: func(c *cli.Context) error {
				args := c.Args()
				var fileURL, msg string
				var err error
				if len(args) == 1 {
					fileURL, msg, err = qn.FetchWithoutKey(args[0])
				}
				if len(args) >= 2 {
					fileURL, msg, err = qn.Fetch(args[0], args[1])
				}
				if err == nil {
					fmt.Println("抓取成功，文件信息为：")
					fmt.Println(fileURL)
					fmt.Println(msg)
				}
				return err
			},
		},
	}
}

func appAction(c *cli.Context) error {
	firstArg := c.Args().First()
	fileList, err := qn.List(firstArg)
	if err != nil {
		return err
	}
	if len(fileList) == 0 {
		fmt.Printf("Prefix %s 无资源\n", firstArg)
		return nil
	}
	for i, file := range fileList {
		fmt.Printf("%4d. %7s %s\n", i+1, utils.GetFileSize(file.Fsize), qn.Origin+file.Key)
	}
	return nil
}
