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
