/**
 *
 * (c) Copyright Ascensio System SIA 2023
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

import {http} from './http';

const ONLYOFFICE_PLUGIN_API = `/plugins/${pluginName}/api`;
export const ONLYOFFICE_PLUGIN_PERMISSIONS = `${ONLYOFFICE_PLUGIN_API}/permissions`;
export const ONLYOFFICE_PLUGIN_GET_CODE = `${ONLYOFFICE_PLUGIN_API}/code`;

export async function get<T>(path: string, config?: RequestInit): Promise<T> {
    const init = {method: 'GET', ...config};
    return http<T>(path, init);
}

export async function post<T, U>(path: string, body: T, config?: RequestInit): Promise<U> {
    const init = {method: 'POST', body: JSON.stringify(body), ...config};
    return http<U>(path, init);
}
