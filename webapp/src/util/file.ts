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
import {FileInfo} from 'mattermost-redux/types/files';

import xlsx from 'public/images/xlsx.svg';
import docx from 'public/images/docx.svg';
import pptx from 'public/images/pptx.svg';
import word from 'public/images/generic_word.svg';
import cell from 'public/images/generic_cell.svg';
import slide from 'public/images/generic_slide.svg';

import {getCookie} from './cookie';

const ONLYOFFICE_CELL = 'cell';
const ONLYOFFICE_WORD = 'word';
const ONLYOFFICE_SLIDE = 'slide';

const EditExtensionsMap = new Map([
    ['docx', ONLYOFFICE_WORD],
    ['xlsx', ONLYOFFICE_CELL],
    ['pptx', ONLYOFFICE_SLIDE],
]);

const AllowedExtensionsMap = new Map([
    ['xls', ONLYOFFICE_CELL],
    ['xlsx', ONLYOFFICE_CELL],
    ['xlsm', ONLYOFFICE_CELL],
    ['xlt', ONLYOFFICE_CELL],
    ['xltx', ONLYOFFICE_CELL],
    ['xltm', ONLYOFFICE_CELL],
    ['ods', ONLYOFFICE_CELL],
    ['fods', ONLYOFFICE_CELL],
    ['ots', ONLYOFFICE_CELL],
    ['csv', ONLYOFFICE_CELL],
    ['pps', ONLYOFFICE_SLIDE],
    ['ppsx', ONLYOFFICE_SLIDE],
    ['ppsm', ONLYOFFICE_SLIDE],
    ['ppt', ONLYOFFICE_SLIDE],
    ['pptx', ONLYOFFICE_SLIDE],
    ['pptm', ONLYOFFICE_SLIDE],
    ['pot', ONLYOFFICE_SLIDE],
    ['potx', ONLYOFFICE_SLIDE],
    ['potm', ONLYOFFICE_SLIDE],
    ['odp', ONLYOFFICE_SLIDE],
    ['fodp', ONLYOFFICE_SLIDE],
    ['otp', ONLYOFFICE_SLIDE],
    ['doc', ONLYOFFICE_WORD],
    ['docx', ONLYOFFICE_WORD],
    ['docm', ONLYOFFICE_WORD],
    ['dot', ONLYOFFICE_WORD],
    ['dotx', ONLYOFFICE_WORD],
    ['dotm', ONLYOFFICE_WORD],
    ['odt', ONLYOFFICE_WORD],
    ['fodt', ONLYOFFICE_WORD],
    ['ott', ONLYOFFICE_WORD],
    ['rtf', ONLYOFFICE_WORD],
]);

const ExtensionIcons = new Map([
    ['xlsx', xlsx],
    ['pptx', pptx],
    ['docx', docx],
    [ONLYOFFICE_WORD, word],
    [ONLYOFFICE_CELL, cell],
    [ONLYOFFICE_SLIDE, slide],
]);

export function getIconByExt(fileExt: string): string {
    const sanitized = fileExt.replaceAll('.', '');
    if (ExtensionIcons.has(sanitized)) {
        return ExtensionIcons.get(sanitized)!;
    }
    return ExtensionIcons.get(getFileTypeByExt(sanitized))!;
}

export function getFileTypeByExt(fileExt: string): string {
    const sanitized = fileExt.replaceAll('.', '');
    if (AllowedExtensionsMap.has(sanitized)) {
        return AllowedExtensionsMap.get(sanitized)!;
    }
    return '';
}

export function isExtensionSupported(fileExt: string, editOnly?: boolean): boolean {
    const sanitized = fileExt.replaceAll('.', '');
    if (editOnly) {
        if (EditExtensionsMap.has(sanitized)) {
            return true;
        }
        return false;
    }
    if (AllowedExtensionsMap.has(sanitized)) {
        return true;
    }

    return false;
}

export function isFileAuthor(fileInfo: FileInfo): boolean {
    const userId: string = getCookie('MMUSERID');

    if (userId) {
        return fileInfo.user_id === userId;
    }

    return false;
}

const fileHelper = {
    getIconByExt,
    getFileTypeByExt,
    isExtensionSupported,
    isFileAuthor,
};

export default fileHelper;
