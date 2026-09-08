// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {render, screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import FormatMultiSelectTable from './FormatMultiSelectTable';
import Formats from './Formats';

describe('FormatMultiSelectTable', () => {
    const options = [
        {label: 'DOCX', value: 'docx'},
        {label: 'XLSX', value: 'xlsx'},
    ];

    it('shows empty state', () => {
        render(
            <FormatMultiSelectTable
                id='formats'
                label='Formats'
                value=''
                disabled={false}
                onChange={jest.fn()}
                setSaveNeeded={jest.fn()}
                options={[]}
            />,
        );

        expect(screen.getByText(/no formats available/i)).toBeInTheDocument();
    });

    it('selects all formats by default when value is empty', async () => {
        render(
            <FormatMultiSelectTable
                id='formats'
                label='Formats'
                value=''
                disabled={false}
                onChange={jest.fn()}
                setSaveNeeded={jest.fn()}
                options={options}
                helpText='help'
            />,
        );

        await waitFor(() => {
            expect(screen.getByLabelText(/deselect all formats/i)).toBeChecked();
        });

        expect(screen.getByText('help')).toBeInTheDocument();
    });

    it('toggles individual formats and select-all', async () => {
        const onChange = jest.fn();
        const setSaveNeeded = jest.fn();

        render(
            <FormatMultiSelectTable
                id='formats'
                label='Formats'
                value='docx, xlsx'
                disabled={false}
                onChange={onChange}
                setSaveNeeded={setSaveNeeded}
                options={options}
            />,
        );

        await waitFor(() => {
            expect(screen.getByLabelText(/docx format/i)).toBeChecked();
        });

        await userEvent.click(screen.getByLabelText(/docx format/i));
        expect(onChange).toHaveBeenCalledWith('formats', 'xlsx');
        expect(setSaveNeeded).toHaveBeenCalled();

        await userEvent.click(screen.getByLabelText(/select all formats|deselect all formats/i));
        expect(onChange).toHaveBeenCalled();
    });

    it('loads none value as empty selection', async () => {
        render(
            <FormatMultiSelectTable
                id='formats'
                label='Formats'
                value='none'
                disabled={false}
                onChange={jest.fn()}
                setSaveNeeded={jest.fn()}
                options={options}
            />,
        );

        await waitFor(() => {
            expect(screen.getByLabelText(/select all formats/i)).not.toBeChecked();
        });
    });
});

describe('Formats', () => {
    it('passes derived options into FormatMultiSelectTable', () => {
        render(
            <Formats
                id='Formats'
                label='Allowed formats'
                value=''
                disabled={false}
                onChange={jest.fn()}
                setSaveNeeded={jest.fn()}
            />,
        );

        expect(screen.getByText('Allowed formats')).toBeInTheDocument();
        expect(screen.getByLabelText(/docx format/i)).toBeInTheDocument();
    });
});
