/* eslint-disable no-console */
import React, {useEffect} from 'react';
import PropTypes from 'prop-types';

import {id as pluginName} from '../../manifest';

// eslint-disable-next-line react/prop-types
const Editor = ({visible, close, theme, fileInfo}) => {
    useEffect(() => {
        if (!visible) {
            return;
        }
        document.getElementById('editorForm').action = `/plugins/${pluginName}/onlyofficeapi/editor`;
        // eslint-disable-next-line react/prop-types
        document.getElementById('file-id').value = fileInfo.id;
        document.getElementById('editorForm').submit();
    }, [fileInfo, visible]);

    if (!visible) {
        return null;
    }

    console.log(fileInfo);

    const style = getStyle(theme);

    return (
        <div
            style={style.backdrop}
            onClick={close}
        >
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
                style={style.modal}
                name='iframeEditor'
            />
        </div>
    );
};

Editor.propTypes = {
    visible: PropTypes.bool.isRequired,
    close: PropTypes.func.isRequired,
    theme: PropTypes.object.isRequired,
};

const getStyle = (theme) => ({
    backdrop: {
        position: 'absolute',
        display: 'flex',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.50)',
        zIndex: 2000,
        alignItems: 'center',
        justifyContent: 'center',
    },
    modal: {
        height: '90%',
        width: '90%',
        padding: '0',
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
});

export default Editor;
