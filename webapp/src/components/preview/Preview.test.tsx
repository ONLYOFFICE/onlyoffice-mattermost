// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import * as fileHelper from 'util/file';

import {render, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {id as pluginId} from 'manifest';
import React from 'react';
import {Provider} from 'react-redux';
import {createStore} from 'redux';

import Preview from './Preview';

jest.mock('util/file', () => ({
    __esModule: true,
    default: {
        getIconByExt: jest.fn(() => 'icon.svg'),
        isExtensionSupported: jest.fn(() => true),
        isFileAuthor: jest.fn(() => true),
    },
}));

describe('Preview', () => {
    const fileInfo = {id: 'f1', name: 'doc.docx', extension: 'docx', user_id: 'u1'} as any;

    const renderPreview = (healthy = true, theme = 'light') => {
        const store = createStore(() => ({
            [`plugins-${pluginId}`]: {health: {healthy}},
        }));
        return render(
            <Provider store={store}>
                <Preview
                    fileInfo={fileInfo}
                    theme={theme}
                    darkTheme={undefined}
                />
            </Provider>,
        );
    };

    it('shows editor and permissions when healthy and author', () => {
        renderPreview(true);
        expect(screen.getByText('doc.docx')).toBeInTheDocument();
        expect(screen.getByAltText('open editor')).toBeInTheDocument();
        expect(screen.getByAltText('permissions button')).toBeInTheDocument();
    });

    it('shows unavailable state when unhealthy', () => {
        renderPreview(false);
        expect(screen.getByTitle(/document server is currently unavailable/i)).toBeInTheDocument();
        expect(screen.queryByAltText('open editor')).not.toBeInTheDocument();
    });

    it('dispatches open editor on click', async () => {
        const dispatch = jest.fn();
        const store = createStore(() => ({
            [`plugins-${pluginId}`]: {health: {healthy: true}},
        }), undefined as any);
        store.dispatch = dispatch;

        render(
            <Provider store={store}>
                <Preview
                    fileInfo={fileInfo}
                    theme='dark'
                    darkTheme='onyx'
                />
            </Provider>,
        );

        await userEvent.click(screen.getByAltText('open editor'));
        expect(dispatch).toHaveBeenCalledWith(expect.objectContaining({
            payload: fileInfo,
        }));
    });

    it('hides permissions when not author', () => {
        (fileHelper.default.isFileAuthor as jest.Mock).mockReturnValueOnce(false);
        renderPreview(true);
        expect(screen.queryByAltText('permissions button')).not.toBeInTheDocument();
    });
});
