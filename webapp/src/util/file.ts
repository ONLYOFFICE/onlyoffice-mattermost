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

import diagram from 'public/images/diagram.svg';
import docx from 'public/images/docx.svg';
import neutral from 'public/images/neutral.svg';
import pdf from 'public/images/pdf.svg';
import pptx from 'public/images/pptx.svg';
import xlsx from 'public/images/xlsx.svg';

import type {FileInfo} from 'mattermost-redux/types/files';

import {getCookie} from './cookie';
import {formatManager, formatHelpers} from './formats';

import type {PluginConfig} from '../api';

let pluginConfig: PluginConfig | null = null;

function isFormatAllowed(extension: string, allowedFormats: string[]): boolean {
    if (!allowedFormats || allowedFormats.length === 0) {
        return false;
    }

    const sanitized = extension.replaceAll('.', '').toLowerCase();
    return allowedFormats.includes(sanitized);
}

export function setPluginConfig(config: PluginConfig): void {
    pluginConfig = config;
}

export function getIconByExt(fileExt: string): string {
    const sanitized = fileExt.replaceAll('.', '');
    const format = formatManager.getFormatByName(sanitized);
    if (format) {
        switch (format.type) {
        case 'word':
            return docx;
        case 'slide':
            return pptx;
        case 'cell':
            return xlsx;
        case 'pdf':
            return pdf;
        case 'diagram':
            return diagram;
        default:
            break;
        }
    }

    return neutral;
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
        if (!formatHelpers.isEditable(format)) {
            return false;
        }
        return pluginConfig ? isFormatAllowed(fileExt, pluginConfig.edit_formats) : false;
    }

    if (!formatHelpers.isViewable(format)) {
        return false;
    }

    return pluginConfig ? isFormatAllowed(fileExt, pluginConfig.view_formats) : false;
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
