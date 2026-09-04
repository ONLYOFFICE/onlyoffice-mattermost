// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

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
