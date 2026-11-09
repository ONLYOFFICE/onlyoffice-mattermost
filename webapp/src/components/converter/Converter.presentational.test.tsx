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

import {render, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import ConverterActions from './ConverterActions';
import ConverterError from './ConverterError';
import ConverterFormatSelection from './ConverterFormatSelection';
import ConverterHeader from './ConverterHeader';
import ConverterInfo from './ConverterInfo';
import ConverterPasswordInput from './ConverterPasswordInput';

describe('converter presentational components', () => {
    it('renders ConverterError only when message exists', () => {
        const {rerender} = render(<ConverterError error=''/>);
        expect(screen.queryByText('error')).not.toBeInTheDocument();

        rerender(<ConverterError error='error'/>);
        expect(screen.getByText('error')).toBeInTheDocument();
    });

    it('renders ConverterInfo content', () => {
        render(<ConverterInfo/>);
        expect(screen.getByRole('heading')).toBeInTheDocument();
    });

    it('calls ConverterHeader close handler', async () => {
        const onClose = jest.fn();

        render(
            <ConverterHeader
                theme='light'
                onClose={onClose}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /close/i}));
        expect(onClose).toHaveBeenCalled();
    });

    it('disables convert when password or format is required', () => {
        const {rerender} = render(
            <ConverterActions
                loading={false}
                needsPassword={true}
                password=''
                needsFormatSelection={false}
                selectedFormat={null}
                onClose={jest.fn()}
                onConvert={jest.fn()}
            />,
        );

        expect(screen.getByRole('button', {name: /convert/i})).toBeDisabled();

        rerender(
            <ConverterActions
                loading={true}
                needsPassword={false}
                password=''
                needsFormatSelection={false}
                selectedFormat={null}
                onClose={jest.fn()}
                onConvert={jest.fn()}
            />,
        );

        expect(screen.getByRole('button', {name: /converting/i})).toBeDisabled();
    });

    it('invokes convert when enabled', async () => {
        const onConvert = jest.fn();

        render(
            <ConverterActions
                loading={false}
                needsPassword={false}
                password=''
                needsFormatSelection={false}
                selectedFormat={null}
                onClose={jest.fn()}
                onConvert={onConvert}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /convert/i}));
        expect(onConvert).toHaveBeenCalled();
    });

    it('updates password input', async () => {
        const onPasswordChange = jest.fn();

        render(
            <ConverterPasswordInput
                password=''
                onPasswordChange={onPasswordChange}
            />,
        );

        await userEvent.type(screen.getByPlaceholderText(/password/i), 'secret');
        expect(onPasswordChange).toHaveBeenCalled();
    });

    it('selects conversion formats', async () => {
        const onFormatSelect = jest.fn();

        render(
            <ConverterFormatSelection
                selectedFormat={null}
                theme='light'
                darkTheme={undefined}
                onFormatSelect={onFormatSelect}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /document/i}));
        await userEvent.click(screen.getByRole('button', {name: /spreadsheet/i}));
        expect(onFormatSelect).toHaveBeenCalledWith('docx');
        expect(onFormatSelect).toHaveBeenCalledWith('xlsx');
    });
});
