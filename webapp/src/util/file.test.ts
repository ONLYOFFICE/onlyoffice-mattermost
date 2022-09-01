/**
 *
 * (c) Copyright Ascensio System SIA 2022
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
import {getIconByExt, getFileTypeByExt, isExtensionSupported} from './file';

describe('ONLYOFFICE File utilities', () => {
    it('get icon by file extension (docx)', () => {
        expect(getIconByExt('docx')).toBeDefined();
    });

    it('get invalid icon by extension (unknown)', () => {
        expect(getIconByExt('unknown')).toBeUndefined();
    });

    it('get file type by file extension (docx)', () => {
        expect(getFileTypeByExt('.docx')).toEqual('word');
    });

    it('get file type by file extension (unknown)', () => {
        expect(getFileTypeByExt('unknown')).toEqual('');
    });

    it('check editable format (docx)', () => {
        expect(isExtensionSupported('docx', true)).toBeTruthy();
    });

    it('check allowed format (docx)', () => {
        expect(isExtensionSupported('docx')).toBeTruthy();
    });

    it('check allowed format (unknown)', () => {
        expect(isExtensionSupported('unknown')).toBeFalsy();
    });
});

