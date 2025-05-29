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

import docx from 'public/images/docx.svg';
import cell from 'public/images/generic_cell.svg';
import slide from 'public/images/generic_slide.svg';
import word from 'public/images/generic_word.svg';
import pptx from 'public/images/pptx.svg';
import xlsx from 'public/images/xlsx.svg';

import type {FileInfo} from 'mattermost-redux/types/files';

import {formatManager, formatHelpers} from './formats';
import {getCookie} from './cookie';

const ExtensionIcons = new Map([
    ['xlsx', xlsx],
    ['pptx', pptx],
    ['docx', docx],
    ['word', word],
    ['cell', cell],
    ['slide', slide],
]);

export function getIconByExt(fileExt: string): string {
    const sanitized = fileExt.replaceAll('.', '');
    if (ExtensionIcons.has(sanitized)) {
        return ExtensionIcons.get(sanitized)!;
    }
    const format = formatManager.getFormatByName(sanitized);
    return format ? ExtensionIcons.get(format.type)! : '';
}

export function getFileTypeByExt(fileExt: string): string {
    const sanitized = fileExt.replaceAll('.', '');
    const format = formatManager.getFormatByName(sanitized);
    return format ? format.type : '';
}

export function isConvertSupported(fileExt: string): boolean {
    const sanitized = fileExt.replaceAll('.', '');
    const format = formatManager.getFormatByName(sanitized);
    return format ? formatHelpers.isAutoConvertable(format) : false;
}

export function isExtensionSupported(fileExt: string, editOnly?: boolean): boolean {
    const sanitized = fileExt.replaceAll('.', '');
    const format = formatManager.getFormatByName(sanitized);
    
    if (!format) {
        return false;
    }

    if (editOnly) {
        return formatHelpers.isEditable(format);
    }

    return formatHelpers.isViewable(format);
}

export function isFileAuthor(fileInfo: FileInfo): boolean {
    const userId: string = getCookie('MMUSERID');
    return userId ? fileInfo.user_id === userId : false;
}

const fileHelper = {
    getIconByExt,
    getFileTypeByExt,
    isExtensionSupported,
    isFileAuthor,
};

export default fileHelper;
