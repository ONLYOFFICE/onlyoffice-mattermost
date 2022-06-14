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

package main

const MATTERMOST_USER_COOKIE = "MMUSERID"
const MATTERMOST_USER_HEADER = "Mattermost-User-Id"
const MATTERMOST_COPY_POST_LINK_SEPARATOR = "/pl/"

const ONLYOFFICE_AUTHORIZATION_USERID_HEADER = "ONLYOFFICE_USER_ID"
const ONLYOFFICE_AUTHORIZATION_USERNAME_HEADER = "ONLYOFFICE_USERNAME"
const ONLYOFFICE_FILEVALIDATION_FILEID_HEADER = "ONLYOFFICE_FILEID"
const ONLYOFFICE_FILEVALIDATION_POSTID_HEADER = "ONLYOFFICE_POSTID"
const ONYLOFFICE_CHANNELVALIDATION_CHANNELID_HEADER = "ONLYOFFICE_CHANNELID"

const ONLYOFFICE_LOGGER_PREFIX = "[ONLYOFFICE]: "
const ONLYOFFICE_BOT_LOGGER_PREFIX = "[ONLYOFFICE_BOT]: "

const ONLYOFFICE_COMMAND_SERVICE string = "coauthoring/CommandService.ashx"
const ONLYOFFICE_API_JS string = "web-apps/apps/api/documents/api.js"
const ONLYOFFICE_CONVERTER string = ""

const ONLYOFFICE_API_PATH = "plugins/com.onlyoffice.mattermost/onlyofficeapi"

const ONLYOFFICE_ROUTE_EDITOR = "/editor"
const ONLYOFFICE_ROUTE_DOWNLOAD = "/download"
const ONLYOFFICE_ROUTE_CALLBACK = "/callback"
const ONLYOFFICE_ROUTE_SET_PERMISSIONS = "/set_file_permissions"
const ONLYOFFICE_ROUTE_GET_PERMISSIONS = "/get_file_permissions"
const ONLYOFFICE_ROUTE_GET_CHANNEL_USERS = "/channel_users"
const ONLYOFFICE_ROUTE_GET_OTP = "/otp"
