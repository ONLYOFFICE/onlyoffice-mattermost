/**
 *
 * (c) Copyright Ascensio System SIA 2021
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

export const apiGET = async (url: string, headers?: HeadersInit) => {
    let json;
    const response = await fetch(url, {
        method: 'GET',
        headers,
    });

    if (response.body) {
        json = await response.json();
    }

    return json;
};

export const apiPOST = async (url: string, body: string, headers?: HeadersInit) => {
    try {
        await fetch(url, {
            method: 'POST',
            headers,
            body,
        });
    } catch {
        throw new Error('API POST call error');
    }
};
