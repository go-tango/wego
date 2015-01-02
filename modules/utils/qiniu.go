package utils

import (
	"github.com/qiniu/api/fop"
	"github.com/qiniu/api/rs"
)

func GetQiniuUptoken(bucketName string) string {
	putPolicy := rs.PutPolicy{
		Scope: bucketName,
	}
	return putPolicy.Token(nil)
}
func GetQiniuPrivateDownloadUrl(domain, key string) string {
	baseUrl := rs.MakeBaseUrl(domain, key)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(baseUrl, nil)
}

func GetQiniuPublicDownloadUrl(domain, key string) string {
	return "http://" + domain + "/" + key
}

func GetQiniuZoomViewUrl(imageUrl string, width, height int) string {
	var view = fop.ImageView{
		Width:   width,
		Height:  height,
		Quality: 100,
	}

	return view.MakeRequest(imageUrl)
}
