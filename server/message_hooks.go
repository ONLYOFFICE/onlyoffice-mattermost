package main

import (
	"io"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// type ONLYOFFICE_POST_WRAPPER struct {
// 	PostId    string
// 	FileInfos []model.FileInfo
// }

func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	return post, ""
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	// customId := post.GetProp("ONLYOFFICE_ID")

	// if customId != nil {
	// 	var fileInfos []model.FileInfo = []model.FileInfo{}

	// 	for _, fildId := range post.FileIds {
	// 		fileInfo, fileInfoErr := p.API.GetFileInfo(fildId)

	// 		if fileInfoErr != nil {
	// 			continue
	// 		}

	// 		if utils.IsExtensionSupported(fileInfo.Extension) {
	// 			fileInfos = append(fileInfos, *fileInfo)
	// 		}
	// 	}

	// 	if len(fileInfos) > 0 {
	// 		wrapper := ONLYOFFICE_POST_WRAPPER{
	// 			PostId:    post.Id,
	// 			FileInfos: fileInfos,
	// 		}

	// 		wrapperBuffer := new(bytes.Buffer)
	// 		json.NewEncoder(wrapperBuffer).Encode(wrapper)

	// 		customId := fmt.Sprintf("%v", customId)
	// 		p.API.KVSetWithExpiry(customId, wrapperBuffer.Bytes(), 60*30)
	// 	}
	// }
}

func (p *Plugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {

	return info, ""
}
