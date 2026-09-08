// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {render, screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import * as api from 'api';
import React from 'react';
import {Provider} from 'react-redux';
import {createStore} from 'redux';

import Manager from './Manager';

jest.mock('api', () => ({
    get: jest.fn(),
    post: jest.fn(),
    ONLYOFFICE_PLUGIN_GET_CODE: '/api/code',
    ONLYOFFICE_PLUGIN_CREATE: '/api/create',
}));

jest.mock('mattermost-redux/selectors/entities/channels', () => ({
    getCurrentChannelId: () => 'channel-1',
}));

describe('Manager', () => {
    const close = jest.fn(() => jest.fn());
    const store = createStore(() => ({}));

    beforeEach(() => {
        jest.clearAllMocks();
    });

    const renderManager = (visible = true) => render(
        <Provider store={store}>
            <Manager
                visible={visible}
                theme='light'
                darkTheme={undefined}
                close={close}
            />
        </Provider>,
    );

    it('renders nothing when hidden', () => {
        const {container} = renderManager(false);
        expect(container).toBeEmptyDOMElement();
    });

    it('validates empty file name', async () => {
        renderManager();

        const input = screen.getByPlaceholderText(/file name/i);
        await userEvent.clear(input);

        expect(await screen.findByText(/cannot be empty/i)).toBeInTheDocument();
    });

    it('creates a file and closes', async () => {
        jest.useFakeTimers();
        const user = userEvent.setup({advanceTimers: jest.advanceTimersByTime});
        (api.get as jest.Mock).mockResolvedValue('code');
        (api.post as jest.Mock).mockResolvedValue({});

        renderManager();

        await user.click(screen.getByRole('button', {name: /create/i}));
        await waitFor(() => expect(api.post).toHaveBeenCalled());

        jest.advanceTimersByTime(300);

        await waitFor(() => expect(close).toHaveBeenCalled());

        jest.useRealTimers();
    });

    it('shows create failure', async () => {
        (api.get as jest.Mock).mockRejectedValue(new Error('fail'));

        renderManager();

        await userEvent.click(screen.getByRole('button', {name: /create/i}));
        expect(await screen.findByText(/failed to create/i)).toBeInTheDocument();
    });
});
