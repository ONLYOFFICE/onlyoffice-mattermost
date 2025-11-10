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

import React, {useState, useEffect, useCallback} from 'react';

import 'public/scss/format.scss';

interface Props {
    id: string;
    label: string;
    helpText?: string;
    value: string;
    disabled: boolean;
    onChange: (id: string, value: string) => void;
    setSaveNeeded: () => void;
    options: FormatOption[];
}

interface HeaderProps {
    id: string;
    checked: boolean;
    disabled: boolean;
    onChange: () => void;
}

interface ItemProps {
    id: string;
    option: FormatOption;
    checked: boolean;
    disabled: boolean;
    onChange: (value: string) => void;
}

interface GridProps {
    id: string;
    label: string;
    options: FormatOption[];
    selectedFormats: string[];
    disabled: boolean;
    onToggle: (value: string) => void;
}

function Header({id, checked, disabled, onChange}: HeaderProps) {
    return (
        <div className='onlyoffice-format-table__header'>
            <label
                className='onlyoffice-format-table__select-all'
                htmlFor={id}
            >
                <input
                    id={id}
                    type='checkbox'
                    checked={checked}
                    onChange={onChange}
                    disabled={disabled}
                    style={{cursor: disabled ? 'not-allowed' : 'pointer'}}
                    aria-label={checked ? 'Deselect all formats' : 'Select all formats'}
                />
                <span>
                    {checked ? 'Deselect All Formats' : 'Select All Formats'}
                </span>
            </label>
        </div>
    );
}

function Item({id, option, checked, disabled, onChange}: ItemProps) {
    const optionId = `${id}-${option.value}`;
    return (
        <label
            className='onlyoffice-format-table__option'
            htmlFor={optionId}
        >
            <input
                id={optionId}
                type='checkbox'
                checked={checked}
                onChange={() => onChange(option.value)}
                disabled={disabled}
                aria-label={`${option.label} format`}
            />
            <span className='onlyoffice-format-table__option-label'>
                {option.label}
            </span>
        </label>
    );
}

function Grid({id, label, options, selectedFormats, disabled, onToggle}: GridProps) {
    return (
        <div
            className='onlyoffice-format-table__grid'
            role='group'
            aria-label={`${label} options`}
        >
            {options.map((option) => (
                <Item
                    key={option.value}
                    id={id}
                    option={option}
                    checked={selectedFormats.includes(option.value)}
                    disabled={disabled}
                    onChange={onToggle}
                />
            ))}
        </div>
    );
}

export default function FormatMultiSelectTable({
    id,
    label,
    helpText,
    value,
    disabled,
    onChange,
    setSaveNeeded,
    options,
}: Props) {
    const [selectedFormats, setSelectedFormats] = useState<string[]>([]);
    const [selectAll, setSelectAll] = useState(false);

    useEffect(() => {
        if (!options || options.length === 0) {
            setSelectedFormats([]);
            setSelectAll(false);
            return;
        }

        if (value && value.trim() !== '' && value.trim().toLowerCase() !== 'none') {
            const formats = value.split(',').map((format) => format.trim());
            setSelectedFormats(formats);
            setSelectAll(formats.length === options.length);
        } else if (value && value.trim().toLowerCase() === 'none') {
            setSelectedFormats([]);
            setSelectAll(false);
        } else {
            const allFormats = options.map((option) => option.value);
            setSelectedFormats(allFormats);
            setSelectAll(true);
        }
    }, [value, options]);

    const handleToggle = useCallback((format: string) => {
        if (disabled) {
            return;
        }

        setSelectedFormats((prevFormats) => {
            const isSelected = prevFormats.includes(format);
            const newFormats = isSelected ?
                prevFormats.filter((f) => f !== format) :
                [...prevFormats, format];

            const isAllSelected = newFormats.length === options.length;
            setSelectAll(isAllSelected);

            let newValue: string;
            if (isAllSelected) {
                newValue = '';
            } else if (newFormats.length === 0) {
                newValue = 'none';
            } else {
                newValue = newFormats.join(', ');
            }

            onChange(id, newValue);
            setSaveNeeded();

            return newFormats;
        });
    }, [disabled, options.length, id, onChange, setSaveNeeded]);

    const handleSelectAll = useCallback(() => {
        if (disabled) {
            return;
        }

        if (selectAll) {
            setSelectedFormats([]);
            setSelectAll(false);
            onChange(id, 'none');
            setSaveNeeded();
        } else {
            const newFormats = options.map((option) => option.value);
            setSelectedFormats(newFormats);
            setSelectAll(true);
            onChange(id, '');
            setSaveNeeded();
        }
    }, [disabled, selectAll, options, id, onChange, setSaveNeeded]);

    if (!options || options.length === 0) {
        return (
            <div className='form-group onlyoffice-format-table'>
                <label className='control-label col-sm-4'>
                    {label}
                </label>
                <div className='col-sm-8'>
                    <p className='help-text'>{'No formats available'}</p>
                </div>
            </div>
        );
    }

    return (
        <div className='form-group onlyoffice-format-table'>
            <label
                className='control-label col-sm-4'
                htmlFor={`${id}-select-all`}
            >
                {label}
            </label>
            <div className='col-sm-8'>
                <div className='onlyoffice-format-table__container'>
                    <Header
                        id={`${id}-select-all`}
                        checked={selectAll}
                        disabled={disabled}
                        onChange={handleSelectAll}
                    />
                    <Grid
                        id={id}
                        label={label}
                        options={options}
                        selectedFormats={selectedFormats}
                        disabled={disabled}
                        onToggle={handleToggle}
                    />
                </div>
                {helpText && (
                    <div className='help-text'>
                        {helpText}
                    </div>
                )}
            </div>
        </div>
    );
}

export interface FormatOption {
    label: string;
    value: string;
}
