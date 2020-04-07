package DbClient

import (
	"github.com/gentwolf-shen/gohelper/hashhelper"
	"github.com/gentwolf-shen/gohelper/timehelper"
)

type Auth struct {
	AppKey string
	AppSecret string
}

func (this *Auth) GetAuthToken(sql string) string {
	date := timehelper.UtcDate()
	sqlMd5 := hashhelper.Md5(sql)

	str := this.AppKey + "|" + date + "|" + sqlMd5
	sign := hashhelper.Md5(str + "|" + this.AppSecret)
	return str + "|" + sign
}

