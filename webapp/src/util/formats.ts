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

import formatsData from 'public/formats/onlyoffice-docs-formats.json';

export interface Format {
    name: string;
    type: string;
    actions: Set<string>;
    convert: Set<string>;
    mime: string[];
}

export interface FormatManager {
    escapeFileName(filename: string): string;
    getFormatByName(name: string): Format | undefined;
    getAllFormats(): Map<string, Format>;
}

export class MapFormatManager implements FormatManager {
    private formats: Map<string, Format>;

    constructor() {
        this.formats = new Map();
        this.initializeFormats();
    }

    private initializeFormats(): void {
        interface RawFormat {
            name: string;
            type: string;
            actions: string[];
            convert: string[];
            mime: string[];
        }

        (formatsData as RawFormat[]).forEach((rawFormat) => {
            if (!rawFormat.actions.includes('view')) {
                return;
            }

            this.formats.set(rawFormat.name, {
                name: rawFormat.name,
                type: rawFormat.type,
                actions: new Set(rawFormat.actions),
                convert: new Set(rawFormat.convert),
                mime: rawFormat.mime,
            });
        });
    }

    public escapeFileName(filename: string): string {
        // eslint-disable-next-line no-useless-escape
        return filename.replace(/[\\/]/g, ':');
    }

    public getFormatByName(name: string): Format | undefined {
        return this.formats.get(name);
    }

    public getAllFormats(): Map<string, Format> {
        return new Map(this.formats);
    }
}

export const formatHelpers = {
    isLossyEditable(format: Format): boolean {
        return format.actions.has('lossy-edit');
    },

    isEditable(format: Format): boolean {
        return format.actions.has('edit');
    },

    isViewable(format: Format): boolean {
        return format.actions.has('view');
    },

    isViewOnly(format: Format): boolean {
        return format.actions.has('view') && format.actions.size === 1;
    },

    isFillable(format: Format): boolean {
        return format.actions.has('fill');
    },

    isAutoConvertable(format: Format): boolean {
        return format.actions.has('auto-convert');
    },

    isOpenXMLConvertable(format: Format): boolean {
        return format.convert.has('docx') ||
               format.convert.has('pptx') ||
               format.convert.has('xlsx');
    },

    getOpenXMLExtension(format: Format): string {
        switch (format.type) {
        case 'cell':
            return 'xlsx';
        case 'slide':
            return 'pptx';
        default:
            return 'docx';
        }
    },
};

export const formatManager = new MapFormatManager();
