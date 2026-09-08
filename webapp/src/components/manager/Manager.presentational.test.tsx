// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {render, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import {ManagerIcon} from './Icon';
import ManagerActions from './ManagerActions';
import ManagerError from './ManagerError';
import ManagerForm from './ManagerForm';
import ManagerHeader from './ManagerHeader';

jest.mock('redux/selectors', () => ({
    getCurrentTheme: jest.fn(() => 'light'),
}));

describe('manager presentational components', () => {
    it('renders ManagerError only with message', () => {
        const {rerender} = render(<ManagerError error=''/>);
        expect(screen.queryByText('err')).not.toBeInTheDocument();

        rerender(<ManagerError error='err'/>);

        expect(screen.getByText('err')).toBeInTheDocument();
    });

    it('disables create when loading or error', async () => {
        const onCreate = jest.fn();
        const {rerender} = render(
            <ManagerActions
                loading={false}
                error='bad'
                onClose={jest.fn()}
                onCreate={onCreate}
            />,
        );

        expect(screen.getByRole('button', {name: /create/i})).toBeDisabled();

        rerender(
            <ManagerActions
                loading={false}
                error=''
                onClose={jest.fn()}
                onCreate={onCreate}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /create/i}));

        expect(onCreate).toHaveBeenCalled();
    });

    it('updates form fields', async () => {
        const onFileNameChange = jest.fn();
        const onFileTypeChange = jest.fn();
        render(
            <ManagerForm
                fileType='docx'
                fileName='Doc'
                loading={false}
                error=''
                theme='light'
                darkTheme={undefined}
                onFileTypeChange={onFileTypeChange}
                onFileNameChange={onFileNameChange}
            />,
        );

        await userEvent.type(screen.getByPlaceholderText(/file name/i), 'x');
        expect(onFileNameChange).toHaveBeenCalled();
        await userEvent.selectOptions(screen.getByRole('combobox'), 'xlsx');
        expect(onFileTypeChange).toHaveBeenCalledWith('xlsx');
    });

    it('calls header close handler', async () => {
        const onClose = jest.fn();
        render(
            <ManagerHeader
                theme='light'
                loading={false}
                onClose={onClose}
            />,
        );

        await userEvent.click(screen.getByRole('button', {name: /close/i}));
        expect(onClose).toHaveBeenCalled();
    });

    it('renders manager icon from theme', () => {
        const store = {getState: () => ({})} as any;
        render(<ManagerIcon store={store}/>);
        expect(screen.getByAltText('open manager')).toBeInTheDocument();
    });
});
