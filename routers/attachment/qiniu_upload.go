// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package attachment

import (
	"github.com/go-tango/wego/setting"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/context"

	"github.com/go-tango/wego/modules/attachment"
	"github.com/go-tango/wego/modules/models"
	"github.com/go-tango/wego/modules/utils"
	"github.com/go-tango/wego/routers/base"

	"github.com/lunny/tango"
)

type QiniuUploadRouter struct {
	base.BaseRouter
}

func (this *QiniuUploadRouter) Post() {
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = &result
		this.ServeJson(this.Data)
	}()

	// check permition
	if !this.User.IsActive {
		return
	}

	// get file object
	file, handler, err := this.Req().FormFile("image")
	if err != nil {
		return
	}
	defer file.Close()

	t := time.Now()

	image := models.Image{}
	image.User = &this.User

	// get mime type
	mime := handler.Header.Get("Content-Type")

	// save and resize image
	if err := attachment.SaveImageToQiniu(&image, file, mime, handler.Filename, t, setting.QiniuPostBucket); err != nil {
		beego.Error(err)
		return
	}

	result["link"] = image.LinkMiddle()
	result["success"] = true

}

func QiniuImage(ctx *tango.Context) {
	var imageName = path.Base(ctx.Req().RequestURI)
	var imageKey string
	var imageSize string
	if i := strings.IndexRune(imageName, '.'); i == -1 {
		return
	} else {
		imageSize = imageName[i+1:]
		if j := strings.IndexRune(imageSize, '.'); j == -1 {
			imageSize = "full"
		} else {
			imageSize = imageSize[:j]
		}
		imageKey = imageName[:i]
	}

	var image = models.Image{
		Token: imageKey,
	}
	err := image.Read("Token")
	if err != nil {
		return
	}
	var imageWidth = image.Width
	var imageHeight = image.Height
	var zoomRatio = setting.ImageSizeMiddle / imageWidth
	if imageWidth > setting.ImageSizeMiddle {
		imageWidth = setting.ImageSizeMiddle
	}
	imageHeight *= zoomRatio

	var imageUrl = utils.GetQiniuPublicDownloadUrl(setting.QiniuPostDomain, imageKey)
	var zoomImageUrl = utils.GetQiniuZoomViewUrl(imageUrl, imageWidth, imageHeight)
	resp, err := http.Get(zoomImageUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ctx.ResponseWriter.Write(body)
}
