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

import React, {useMemo} from 'react';

import FormatMultiSelectTable from './FormatMultiSelectTable';
import type {FormatOption} from './FormatMultiSelectTable';

import formatsData from 'public/formats/onlyoffice-docs-formats.json';

interface Props {
    id: string;
    label: string;
    value: string;
    disabled: boolean;
    onChange: (id: string, value: string) => void;
    setSaveNeeded: () => void;
}

export default function ViewFormats({
    id,
    label,
    value,
    disabled,
    onChange,
    setSaveNeeded,
}: Props) {
    const viewFormats = useMemo(() => {
        const formats: FormatOption[] = [];
        formatsData.forEach((format: any) => {
            if (format.actions && format.actions.length > 0) {
                const hasViewAction = format.actions.some((action: string) => 
                    action === 'view' || 
                    action === 'lossy-edit' || 
                    action === 'auto-convert'
                );
                
                if (hasViewAction && format.name) {
                    formats.push({
                        label: format.name.toUpperCase(),
                        value: format.name.toLowerCase(),
                    });
                }
            }
        });
        
        return formats.sort((a, b) => a.label.localeCompare(b.label));
    }, []);

    return (
        <FormatMultiSelectTable
            id={id}
            label={label}
            value={value}
            disabled={disabled}
            onChange={onChange}
            setSaveNeeded={setSaveNeeded}
            options={viewFormats}
            helpText="Select file formats that are allowed for viewing in ONLYOFFICE. All formats are enabled by default. Uncheck formats to disable them."
        />
    );
}

