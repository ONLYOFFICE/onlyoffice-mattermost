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

import 'isomorphic-fetch';
import '@testing-library/jest-dom';

const localStorageStore = {};

const localStorageMock = {
    getItem: (key) => (Object.prototype.hasOwnProperty.call(localStorageStore, key) ? localStorageStore[key] : null),
    setItem: (key, value) => {
        localStorageStore[key] = String(value);
    },
    removeItem: (key) => {
        delete localStorageStore[key];
    },
    clear: () => {
        Object.keys(localStorageStore).forEach((key) => delete localStorageStore[key]);
    },
};

Object.defineProperty(window, 'localStorage', {
    value: localStorageMock,
    configurable: true,
});

Object.defineProperty(global, 'localStorage', {
    value: localStorageMock,
    configurable: true,
});
