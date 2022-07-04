/**
 *
 * (c) Copyright Ascensio System SIA 2022
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package bot

import "github.com/mattermost/mattermost-server/v6/plugin"

type Bot interface {
	BotCreateDM(message string, userID string)
	BotCreatePost(message string, channelID string)
	BotCreateReply(message string, channelID string, parentID string)
}

func NewBot(ID string, API plugin.API) Bot {
	return onlyofficeBot{
		ID:  ID,
		API: API,
	}
}
