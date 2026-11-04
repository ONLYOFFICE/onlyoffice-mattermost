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

import {getTranslations} from 'util/lang';

import type {FileDetails, FileDetailsUser} from 'api';
import {getFileDetails} from 'api';
import React, {useEffect, useState} from 'react';

type Props = {
    fileId: string;
    theme: string;
};

export function formatSize(bytes: number): string {
    if (bytes < 1024) {
        return `${bytes} B`;
    }
    const units = ['KB', 'MB', 'GB'];
    let value = bytes / 1024;
    let unit = 0;
    while (value >= 1024 && unit < units.length - 1) {
        value /= 1024;
        unit++;
    }
    return `${value.toFixed(value < 10 ? 1 : 0)} ${units[unit]}`;
}

function formatDate(ms: number): string {
    const locale = window.localStorage.getItem('onlyoffice_locale') || undefined;
    return new Date(ms).toLocaleString(locale, {dateStyle: 'medium', timeStyle: 'short'});
}

function userName(user?: FileDetailsUser): string {
    if (!user) {
        return '';
    }
    return user.username ? `@${user.username}` : user.display_name || user.id;
}

// Extra information about the file shown in the preview: who uploaded it and
// when, the latest save made in ONLYOFFICE, the actual size and what the
// current user may do with the file.
export default function FileDetailsRows({fileId, theme}: Props) {
    const i18n = getTranslations() as Record<string, string>;
    const [details, setDetails] = useState<FileDetails | null>(null);

    useEffect(() => {
        let cancelled = false;
        const load = () => {
            getFileDetails(fileId).
                then((d) => !cancelled && setDetails(d)).
                catch(() => !cancelled && setDetails(null));
        };
        load();

        // The document is edited in another tab: refresh when the user comes back.
        const onVisible = () => {
            if (document.visibilityState === 'visible') {
                load();
            }
        };
        window.addEventListener('focus', load);
        document.addEventListener('visibilitychange', onVisible);
        return () => {
            cancelled = true;
            window.removeEventListener('focus', load);
            document.removeEventListener('visibilitychange', onVisible);
        };
    }, [fileId]);

    if (!details) {
        return null;
    }

    const rows: Array<[string, string]> = [];
    rows.push([i18n['preview.details.uploaded'], [userName(details.uploaded_by), formatDate(details.uploaded_at)].filter(Boolean).join(', ')]);
    if (details.edited_at) {
        rows.push([i18n['preview.details.edited'], [userName(details.edited_by), formatDate(details.edited_at)].filter(Boolean).join(', ')]);
    }
    rows.push([i18n['preview.details.size'], formatSize(details.size)]);
    let access = details.can_edit ? i18n['preview.details.access_edit'] : i18n['preview.details.access_view'];
    if (details.is_owner) {
        access += ` (${i18n['preview.details.owner']})`;
    }
    rows.push([i18n['preview.details.access'], access]);

    return (
        <table
            className='onlyoffice_preview__details'
            data-theme={theme}
        >
            <tbody>
                {rows.map(([label, value]) => (
                    <tr key={label}>
                        <td className='onlyoffice_preview__details-label'>{label}</td>
                        <td className='onlyoffice_preview__details-value'>{value}</td>
                    </tr>
                ))}
            </tbody>
        </table>
    );
}
