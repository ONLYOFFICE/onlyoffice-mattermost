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

import {ONLYOFFICE_ERROR_EVENT, ONLYOFFICE_READY_EVENT} from 'util/const';

import {render, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import Editor from './Editor';
import EditorLoader from './EditorLoader';

describe('EditorLoader', () => {
    it('hides after ready event', () => {
        const {container} = render(<EditorLoader theme='light'/>);
        expect(screen.getByRole('button', {name: /close/i})).toBeInTheDocument();

        window.dispatchEvent(new Event(ONLYOFFICE_READY_EVENT));
        expect(container).toBeEmptyDOMElement();
    });

    it('shows error message from event detail', async () => {
        render(<EditorLoader theme='light'/>);

        window.dispatchEvent(new CustomEvent(ONLYOFFICE_ERROR_EVENT, {
            detail: {messageKey: 'editor.events.unauthorized'},
        }));

        expect(await screen.findByText(/unauthorized/i)).toBeInTheDocument();
    });

    it('dispatches close event from button', async () => {
        const spy = jest.fn();

        window.addEventListener('onlyofficecloseeditor', spy);
        render(<EditorLoader theme='light'/>);

        await userEvent.click(screen.getByRole('button', {name: /close/i}));

        expect(spy).toHaveBeenCalled();

        window.removeEventListener('onlyofficecloseeditor', spy);
    });
});

describe('Editor', () => {
    beforeEach(() => {
        jest.useFakeTimers();
    });

    afterEach(() => {
        jest.useRealTimers();
    });

    it('renders nothing when hidden', () => {
        const {container} = render(
            <Editor
                visible={false}
                theme='light'
                close={jest.fn(() => jest.fn())}
                fileInfo={{id: 'f1'} as any}
            />,
        );

        expect(container).toBeEmptyDOMElement();
    });

    it('renders iframe portal when visible', () => {
        render(
            <Editor
                visible={true}
                theme='dark'
                close={jest.fn(() => jest.fn())}
                fileInfo={{id: 'f1'} as any}
            />,
        );

        expect(document.getElementById('editor-backdrop')).toBeInTheDocument();
        expect(document.querySelector('iframe')).toHaveAttribute(
            'src',
            expect.stringContaining('file=f1'),
        );
    });

    it('closes on ONLYOFFICE close event', () => {
        const close = jest.fn(() => jest.fn());

        render(
            <Editor
                visible={true}
                theme='light'
                close={close}
                fileInfo={{id: 'f1'} as any}
            />,
        );

        window.dispatchEvent(new Event('onlyofficecloseeditor'));

        jest.advanceTimersByTime(280);

        expect(close).toHaveBeenCalled();
    });
});
