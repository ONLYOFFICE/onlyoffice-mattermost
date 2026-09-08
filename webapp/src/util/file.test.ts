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

import {
    getFileTypeByExt,
    getIconByExt,
    isConvertSupported,
    isExtensionSupported,
    isFileAuthor,
    setPluginConfig,
} from './file';

describe('file utils', () => {
    beforeEach(() => {
        setPluginConfig({formats: ['docx', 'xlsx', 'pptx', 'pdf']} as any);
        (document as any).cookie = '';
    });

    it('returns icons and types for known extensions', () => {
        expect(getIconByExt('.docx')).toBeTruthy();
        expect(getIconByExt('.pptx')).toBeTruthy();
        expect(getIconByExt('.xlsx')).toBeTruthy();
        expect(getIconByExt('.pdf')).toBeTruthy();
        expect(getIconByExt('unknown-ext')).toBeTruthy();
        expect(getFileTypeByExt('docx')).toBe('word');
        expect(getFileTypeByExt('xlsx')).toBe('cell');
        expect(getFileTypeByExt('pptx')).toBe('slide');
        expect(getFileTypeByExt('pdf')).toBe('pdf');
        expect(getFileTypeByExt('nope')).toBe('');
    });

    it('checks convert and extension support against plugin config', () => {
        expect(isExtensionSupported('docx')).toBe(true);
        expect(isExtensionSupported('docx', true)).toBe(true);
        expect(isExtensionSupported('not-a-format')).toBe(false);
        expect(isExtensionSupported('not-a-format', true)).toBe(false);

        setPluginConfig(null as any);
        expect(isExtensionSupported('docx')).toBe(false);
        expect(isExtensionSupported('docx', true)).toBe(false);

        setPluginConfig({formats: []} as any);
        expect(isExtensionSupported('docx')).toBe(false);

        setPluginConfig({formats: ['docx']} as any);
        expect(typeof isConvertSupported('doc')).toBe('boolean');
        expect(isConvertSupported('not-a-format')).toBe(false);
    });

    it('detects file authorship from MMUSERID cookie', () => {
        expect(isFileAuthor({user_id: 'u1'} as any)).toBe(false);

        document.cookie = 'MMUSERID=u1';

        expect(isFileAuthor({user_id: 'u1'} as any)).toBe(true);
        expect(isFileAuthor({user_id: 'u2'} as any)).toBe(false);
    });
});
