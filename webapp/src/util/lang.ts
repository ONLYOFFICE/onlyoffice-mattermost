// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/**
 *
 * (c) Copyright Ascensio System SIA 2025
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

import de from 'i18n/de.json';
import en from 'i18n/en.json';
import es from 'i18n/es.json';
import fr from 'i18n/fr.json';
import it from 'i18n/it.json';
import ru from 'i18n/ru.json';

export function getTranslations(locale?: string) {
    if (locale) {
        window.localStorage.setItem('onlyoffice_locale', locale);
    }
    const currentLocale = locale || window.localStorage.getItem('onlyoffice_locale') || 'en';
    switch (currentLocale) {
    case 'de':
        return de;
    case 'en':
        return en;
    case 'es':
        return es;
    case 'fr':
        return fr;
    case 'it':
        return it;
    case 'ru':
        return ru;
    default:
        return en;
    }
}
