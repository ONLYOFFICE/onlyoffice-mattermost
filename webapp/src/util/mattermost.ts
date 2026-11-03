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

// Small replacements for the few mattermost-redux helpers the plugin used.
// mattermost-redux is no longer published for plugins and pulled an outdated
// dependency tree (including a git dependency) into the bundle.

import type {Channel} from '@mattermost/types/channels';
import type {GlobalState} from '@mattermost/types/store';
import type {Team} from '@mattermost/types/teams';
import type {UserProfile} from '@mattermost/types/users';

async function mmFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
    const response = await fetch(path, {
        credentials: 'include',
        ...init,
        headers: {

            // Required by Mattermost for cookie-authenticated API calls (CSRF protection).
            'X-Requested-With': 'XMLHttpRequest',
            'Content-Type': 'application/json',
            ...(init.headers || {}),
        },
    });

    if (!response.ok) {
        throw new Error(`${response.status} ${response.statusText}`);
    }

    return response.json();
}

export function getTeamByName(name: string): Promise<Team> {
    return mmFetch<Team>(`/api/v4/teams/name/${encodeURIComponent(name)}`);
}

export function getChannelByName(teamId: string, name: string): Promise<Channel> {
    return mmFetch<Channel>(`/api/v4/teams/${encodeURIComponent(teamId)}/channels/name/${encodeURIComponent(name)}?include_deleted=true`);
}

export function searchUsers(term: string, options: {in_channel_id?: string; team_id?: string}): Promise<UserProfile[]> {
    return mmFetch<UserProfile[]>('/api/v4/users/search', {
        method: 'POST',
        body: JSON.stringify({term, ...options}),
    });
}

export function getCurrentChannelId(state: GlobalState): string {
    return state.entities.channels.currentChannelId;
}

export type ThemeInfo = {type?: string};

// Mirrors mattermost-redux getTheme(): team-specific theme preference, then the
// global one, then the server default theme.
export function getTheme(state: GlobalState): ThemeInfo {
    const prefs = state.entities.preferences?.myPreferences || {};
    const teamId = state.entities.teams?.currentTeamId || '';
    const pref = prefs[`theme--${teamId}`] || prefs['theme--'];

    if (pref?.value) {
        try {
            return JSON.parse(pref.value) as ThemeInfo;
        } catch {
            // fall through to the default theme
        }
    }

    const defaultTheme = (state.entities.general?.config as Record<string, string> | undefined)?.DefaultTheme;
    return {type: defaultTheme || 'denim'};
}

export function isDarkTheme(theme: ThemeInfo): boolean {
    const type = (theme.type || '').toLowerCase();
    return type === 'onyx' || type === 'indigo';
}
