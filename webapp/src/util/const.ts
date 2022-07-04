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
import {id as pluginName} from 'manifest';

export const ONLYOFFICE_WILDCARD_USER = '*';

export const ONLYOFFICE_PLUGIN_API = `/plugins/${pluginName}/api`;

export const ONLYOFFICE_CLOSE_EVENT = 'onlyofficecloseeditor';
export const ONLYOFFICE_READY_EVENT = 'onlyofficeready';
export const ONLYOFFICE_ERROR_EVENT = 'onlyofficeerror';
