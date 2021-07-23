/* eslint-disable no-console */
/* eslint-disable react/jsx-no-literals */
/* eslint-disable react/prop-types */
import React, {useState} from 'react';
import PropTypes from 'prop-types';

//import {id as pluginName} from '../../manifest';

const PermissionsForm = ({fileId, username, setUsername, permissions, setPermissions}) => {
    const handleOnChange = (permissionName) => {
        const updatedCheckedState = permissions.map((item) =>
            (item.name === permissionName ? Object.assign(item, {checked: !item.checked}) : item),
        );

        setPermissions(updatedCheckedState);
    };

    return (
        <>
            <input
                type='text'
                value={username}
                onChange={(e) => setUsername(e.target.value)}
            />
            <table>
                <thead>
                    <tr>
                        <td>Name</td>
                        <td>Selected</td>
                    </tr>
                </thead>
                <tbody>
                    {permissions.map((permission, index) => {
                        return (
                            <tr key={index}>
                                <td>
                                    {permission.name}
                                </td>
                                <td>
                                    <input
                                        type='checkbox'
                                        id={`${permission.name}-${index}`}
                                        name={permission.name}
                                        value={permission.name}
                                        checked={permissions[index].checked}
                                        onChange={() => handleOnChange(permission.name)}
                                    />
                                </td>
                            </tr>
                        );
                    })}
                </tbody>
            </table>
            <button
                onClick={() => {
                    const body = {};
                    permissions.forEach((permission) => {
                        body[permission.name] = permission.checked;
                    });
                    const OPTIONS = {
                        method: 'POST',
                        body: JSON.stringify(body),
                        headers: {
                            'Content-Type': 'application/json',
                        },
                    };

                    fetch(`http://46.101.101.37:8065/plugins/com.onlyoffice.mattermost-plugin/onlyofficeapi/permissions?fileId=${fileId}&username=${username}`, OPTIONS).catch();
                }}
            >OK</button>
        </>
    );
};

const Permissions = ({visible, close, theme, fileInfo}) => {
    const [permissions, setPermissions] = useState(
        [
            {
                name: 'edit',
                checked: false,
            },
            {
                name: 'download',
                checked: false,
            },
            {
                name: 'print',
                checked: false,
            },
            {
                name: 'comment',
                checked: false,
            },
            {
                name: 'review',
                checked: false,
            },
            {
                name: 'copy',
                checked: false,
            },
        ],
    );
    const [username, setUsername] = useState('*');

    if (!visible) {
        return null;
    }

    const style = getStyle(theme);

    return (
        <div
            style={style.backdrop}
            onClick={close}
        >
            <div
                style={style.modal}
                onClick={(e) => {
                    e.stopPropagation();
                }}
            >
                <PermissionsForm
                    username={username}
                    setUsername={setUsername}
                    permissions={permissions}
                    setPermissions={setPermissions}
                    fileId={fileInfo.id}
                />
            </div>
        </div>
    );
};

Permissions.propTypes = {
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
        height: '250px',
        width: '400px',
        padding: '1em',
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
});

export default Permissions;
