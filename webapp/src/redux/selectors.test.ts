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

import {getTheme} from 'mattermost-redux/selectors/entities/preferences';

import {
    converterModalFileInfo,
    converterModalVisible,
    editorModalFileInfo,
    editorModalVisible,
    getCurrentDarkTheme,
    getCurrentTheme,
    managerModalVisible,
    permissionsModalFileInfo,
    permissionsModalVisible,
} from './selectors';

import {id as pluginId} from '../manifest';

jest.mock('mattermost-redux/selectors/entities/preferences', () => ({
    getTheme: jest.fn(),
}));

describe('redux selectors', () => {
    const fileInfo = {id: 'f1'};
    const state = {
        [`plugins-${pluginId}`]: {
            converterModal: {isVisible: true, fileInfo},
            editorModal: {isVisible: false, fileInfo: null},
            permissionsModal: {isVisible: true, fileInfo},
            managerModal: {isVisible: true},
        },
    } as any;

    it('reads plugin modal state', () => {
        expect(converterModalVisible(state)).toBe(true);
        expect(converterModalFileInfo(state)).toBe(fileInfo);
        expect(editorModalVisible(state)).toBe(false);
        expect(editorModalFileInfo(state)).toBeNull();
        expect(permissionsModalVisible(state)).toBe(true);
        expect(permissionsModalFileInfo(state)).toBe(fileInfo);
        expect(managerModalVisible(state)).toBe(true);
    });

    it('maps theme types to light/dark', () => {
        (getTheme as jest.Mock).mockReturnValue({type: 'default'});
        expect(getCurrentTheme(state)).toBe('light');
        expect(getCurrentDarkTheme(state)).toBeUndefined();

        (getTheme as jest.Mock).mockReturnValue({type: 'indigo'});
        expect(getCurrentTheme(state)).toBe('dark');
        expect(getCurrentDarkTheme(state)).toBe('indigo');

        (getTheme as jest.Mock).mockReturnValue({type: 'onyx'});
        expect(getCurrentTheme(state)).toBe('dark');
        expect(getCurrentDarkTheme(state)).toBe('onyx');
    });
});
