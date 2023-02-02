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
/* eslint-disable max-nested-callbacks */
/* eslint-disable @typescript-eslint/no-explicit-any */
import {Client4} from 'mattermost-redux/client';
import {Channel} from 'mattermost-redux/types/channels';
import {FileInfo} from 'mattermost-redux/types/files';

import {getUniqueMattermostUsers, MattermostUser} from './user';

export function debounce(cb: (...args: any[]) => void, delay: number) {
    let timeout: NodeJS.Timeout;
    return (...args: any[]) => {
        clearTimeout(timeout);
        timeout = setTimeout(() => {
            cb(...args);
        }, delay || 100);
    };
}

export function debounceUsersLoad(channel: Channel | undefined, fileInfo: FileInfo, users: MattermostUser[]) {
    return debounce((input: any, callback: any) => {
        if (!input) {
            return;
        }

        if (channel) {
            (async () => {
                let res = await Client4.searchUsers(input, {
                    in_channel_id: channel.id,
                    team_id: channel.team_id,
                });

                res = res.filter((user) => user.id !== fileInfo.user_id);
                const permissions = getUniqueMattermostUsers(res, users);
                callback(permissions);
            })();
        }
    }, 2000);
}

export const pipe = <T>(...fns: Array<(arg: T) => T>) => (value: T) => fns.reduce((acc, fn) => fn(acc), value);
