import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import {id as pluginName} from '../manifest';

//http://46.101.101.37/web-apps/apps/api/documents/api.js

const Root = ({ visible, close, theme, fileInfo }) => {
    if (!visible) {
        return null;
    }
    useEffect(() => {
        document.getElementById("editorForm").action = `/plugins/${pluginName}/onlyofficeapi/editor`;
        document.getElementById("file-id").value = fileInfo.id;
        document.getElementById("editorForm").submit();
    }, [fileInfo]);
    const style = getStyle(theme);
    return (
        <div
            style={style.backdrop}
            onClick={close}
        >
            <form action="" method="POST" target="iframeEditor" id="editorForm">
                <input id='file-id' name="fileid" value='' type='hidden' />
            </form>
            <iframe style={style.modal} name="iframeEditor" />
        </div>
    );
};

Root.propTypes = {
    visible: PropTypes.bool.isRequired,
    close: PropTypes.func.isRequired,
    theme: PropTypes.object.isRequired,
    subMenu: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
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

export default Root;
