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

import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';
import {openEditor} from 'redux/actions';

import type {FileInfo} from '@mattermost/types/files';

type WebappHistory = {replace: (path: string) => void};

function navigate(path: string): void {
    const history = (window as unknown as {WebappUtils?: {browserHistory?: WebappHistory}}).WebappUtils?.browserHistory;
    if (history) {
        history.replace(path);
    } else {
        window.location.replace(path);
    }
}

export default function CommentRedirect() {
    const dispatch = useDispatch();
    const [failed, setFailed] = useState(false);

    useEffect(() => {
        const params = new URLSearchParams(window.location.search);
        const fileId = params.get('file');
        const action = params.get('action') || undefined;
        const team = window.location.pathname.split('/')[1];
        if (!fileId) {
            setFailed(true);
            return;
        }

        fetch(`/api/v4/files/${encodeURIComponent(fileId)}/info`, {
            credentials: 'include',
            headers: {'X-Requested-With': 'XMLHttpRequest'},
        }).
            then((response) => (response.ok ? response.json() : Promise.reject(new Error(response.statusText)))).
            then((fileInfo: FileInfo) => {
                navigate(`/${team}/pl/${fileInfo.post_id}`);
                openEditor(fileInfo, action)(dispatch);
            }).
            catch(() => setFailed(true));
    }, [dispatch]);

    if (!failed) {
        return null;
    }

    return (
        <div className='onlyoffice-comment-redirect'>
            {'The document is not available.'}
        </div>
    );
}
