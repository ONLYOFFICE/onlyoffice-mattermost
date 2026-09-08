// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {MapFormatManager, formatHelpers, formatManager} from './formats';

describe('formats utils', () => {
    it('loads viewable formats from data', () => {
        expect(formatManager.getAllFormats().size).toBeGreaterThan(0);
        expect(formatManager.getFormatByName('docx')).toBeDefined();
        expect(formatManager.getFormatByName('docx')?.type).toBe('word');
    });

    it('escapes path separators in filenames', () => {
        const manager = new MapFormatManager();
        expect(manager.escapeFileName('a/b\\c.docx')).toBe('a:b:c.docx');
    });

    it('classifies format capabilities', () => {
        const format = {
            name: 'docx',
            type: 'word',
            actions: new Set(['view', 'edit', 'lossy-edit', 'fill', 'auto-convert']),
            convert: new Set(['pdf', 'docx']),
            mime: [],
        };

        expect(formatHelpers.isLossyEditable(format)).toBe(true);
        expect(formatHelpers.isEditable(format)).toBe(true);
        expect(formatHelpers.isViewable(format)).toBe(true);
        expect(formatHelpers.isViewOnly(format)).toBe(false);
        expect(formatHelpers.isFillable(format)).toBe(true);
        expect(formatHelpers.isAutoConvertable(format)).toBe(true);
        expect(formatHelpers.isOpenXMLConvertable(format)).toBe(true);
        expect(formatHelpers.getOpenXMLExtension(format)).toBe('docx');
        expect(formatHelpers.getOpenXMLExtension({...format, type: 'cell'})).toBe('xlsx');
        expect(formatHelpers.getOpenXMLExtension({...format, type: 'slide'})).toBe('pptx');
        expect(formatHelpers.getOpenXMLExtension({...format, type: 'diagram'})).toBe('vsdx');
    });

    it('detects view-only formats', () => {
        const viewOnly = {
            name: 'pdf',
            type: 'pdf',
            actions: new Set(['view']),
            convert: new Set<string>(),
            mime: [],
        };

        expect(formatHelpers.isViewOnly(viewOnly)).toBe(true);
        expect(formatHelpers.isOpenXMLConvertable(viewOnly)).toBe(false);
    });
});
