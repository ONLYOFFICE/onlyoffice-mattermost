// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/**
 *
 * (c) Copyright Ascensio System SIA 2026
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

import {ONLYOFFICE_PLUGIN_API} from 'util/const';

import type {FileInfo} from '@mattermost/types/files';

export function getEditorUrl(fileInfo: FileInfo, theme: string): string {
    const lang = localStorage.getItem('onlyoffice_locale') || 'en';
    return `${ONLYOFFICE_PLUGIN_API}/editor?file=${encodeURIComponent(fileInfo.id)}&lang=${encodeURIComponent(lang)}&dark=${theme === 'dark'}`;
}

// Opens the editor in a separate browser tab/window. One window per file:
// clicking the same file again focuses the already opened editor instead of
// loading a second copy. The opener reference is kept on purpose: editor.html
// uses window.opener to close the tab when the user presses "Close".
export function openEditorInNewWindow(fileInfo: FileInfo, theme: string): void {
    const url = getEditorUrl(fileInfo, theme);
    const win = window.open('', `onlyoffice-editor-${fileInfo.id}`);
    if (!win) {
        // Pop-up blocked: fall back to a plain new tab.
        window.open(url, '_blank');
        return;
    }

    try {
        if (win.location.pathname.endsWith('/api/editor')) {
            win.focus(); // already opened, do not reload the editor
            return;
        }
    } catch {
        // not readable yet: navigate below
    }

    win.location.href = url;
    win.focus();
}
