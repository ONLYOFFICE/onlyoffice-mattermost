// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {getTranslations} from './lang';

describe('lang utils', () => {
    beforeEach(() => {
        window.localStorage.clear();
    });

    it('defaults to english and persists locale', () => {
        const en = getTranslations();
        expect(en).toBeDefined();
        expect(typeof en).toBe('object');

        getTranslations('de');
        expect(window.localStorage.getItem('onlyoffice_locale')).toBe('de');
        expect(getTranslations()).toBe(getTranslations('de'));
    });

    it('returns english for unknown locales', () => {
        expect(getTranslations('zz')).toEqual(getTranslations('en'));
    });

    it('supports known locales', () => {
        ['de', 'en', 'es', 'fr', 'it', 'ru'].forEach((locale) => {
            expect(getTranslations(locale)).toBeDefined();
        });
    });
});
