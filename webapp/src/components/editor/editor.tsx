import React, {useCallback} from 'react';

import {FileInfo} from 'mattermost-redux/types/files';

import {Dispatch} from 'redux';

import {id as pluginName} from 'manifest';

import {EditorLoader} from './editor_loader';

interface EditorProps {
    visible: boolean,
    close: () => (dispatch: Dispatch) => void,
    fileInfo?: FileInfo,
}

const Editor = ({visible, close, fileInfo}: EditorProps) => {
    const handleClose = useCallback(() => {
        if (!visible) {
            return;
        }
        const editorBackdrop = document.getElementById('editor-backdrop');

        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        editorBackdrop!.classList.add('plugin-modal__backdrop_hide');

        setTimeout(() => close(), 300);
    }, [close, visible]);

    const escFunction = useCallback((event: any) => {
        if (event.keyCode === 27) {
            handleClose();
        }
    }, [handleClose]);

    React.useEffect(() => {
        if (!visible || !fileInfo) {
            return;
        }
        (document.getElementById('editorForm') as HTMLFormElement).action = `/plugins/${pluginName}/onlyofficeapi/editor`;
        (document.getElementById('file-id') as HTMLInputElement).value = fileInfo.id;
        (document.getElementById('editorForm') as HTMLFormElement).submit();
        window.addEventListener('ONLYOFFICE_CLOSED', handleClose);
        document.addEventListener('keydown', escFunction, false);

        // eslint-disable-next-line consistent-return
        return () => {
            window.removeEventListener('ONLYOFFICE_CLOSED', handleClose);
            document.removeEventListener('keydown', escFunction, false);
        };
    }, [fileInfo, visible, handleClose, escFunction]);

    return (
        <>
            {visible && (
                <div
                    id='editor-backdrop'
                    className='plugin-modal__backdrop'
                >
                    <EditorLoader/>
                    <form
                        action=''
                        method='POST'
                        target='iframeEditor'
                        id='editorForm'
                    >
                        <input
                            id='file-id'
                            name='fileid'
                            value=''
                            type='hidden'
                        />
                    </form>
                    <iframe
                        className='plugin-modal__frame'
                        name='iframeEditor'
                    />
                </div>
            )}
        </>
    );
};

export default Editor;
