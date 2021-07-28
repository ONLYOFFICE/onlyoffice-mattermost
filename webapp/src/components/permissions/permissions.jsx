/* eslint-disable no-unused-vars */
/* eslint-disable no-console */
/* eslint-disable react/jsx-no-literals */
/* eslint-disable react/prop-types */
import axios from 'axios';
import React, {useState} from 'react';
import PropTypes from 'prop-types';
import {SelectField, Button} from 'evergreen-ui';

import {UsersTable} from './users_table';

//import {id as pluginName} from '../../manifest';

let PERMISSIONS_TO_SEND = [];

const READ_ONLY = {
    edit: false,
};

const FULL_ACCESS = {
    edit: true,
};

const Permissions = ({visible, close, theme, fileInfo}) => {
    const [options] = useState(['Read only', 'Full access']);
    const [selectedOption, setSelectedOption] = useState(options[0]);
    const [selected, setSelected] = useState([]);
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
            <div
                style={style.modal}
                onClick={(e) => {
                    e.stopPropagation();
                }}
            >
                <div style={{display: 'flex', justifyContent: 'space-evenly', alignItems: 'center', width: '100%', height: '85%'}}>
                    <SelectField
                        label=''
                        value={selectedOption}
                        onChange={(e) => setSelectedOption(e.target.value)}
                    >
                        <option
                            value={options[0]}
                            selected={true}
                        >
                            {options[0]}
                        </option>
                        {options.filter((_, index) => index !== 0).map((option) => {
                            return (
                                <option
                                    key={option}
                                    value={option}
                                >
                                    {option}
                                </option>
                            );
                        })}
                    </SelectField>
                    <UsersTable
                        fileInfoId={fileInfo.id}
                        selected={selected}
                        setSelected={setSelected}
                    />
                </div>
                <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center', width: '100%', height: '12%', overflow: 'hidden', marginTop: '1rem'}}>
                    <Button
                        marginRight={50}
                        style={{width: '15%'}}
                        appearance='primary'
                        intent='success'
                        disabled={selected.length === 0}
                        onClick={() => {
                            if (selected === '*') {
                                PERMISSIONS_TO_SEND.push({
                                    Username: '*',
                                    FileId: fileInfo.id,
                                    Permissions: selectedOption === options[0] ? READ_ONLY : FULL_ACCESS,
                                });
                            } else {
                                selected.forEach((user) => {
                                    if (selectedOption === options[0]) {
                                        PERMISSIONS_TO_SEND.push({
                                            Username: user,
                                            FileId: fileInfo.id,
                                            Permissions: READ_ONLY,
                                        });
                                    } else if (selectedOption === options[1]) {
                                        PERMISSIONS_TO_SEND.push({
                                            Username: user,
                                            FileId: fileInfo.id,
                                            Permissions: FULL_ACCESS,
                                        });
                                    }
                                });
                            }

                            axios({
                                url: `http://46.101.101.37:8065/plugins/com.onlyoffice.mattermost-plugin/onlyofficeapi/file_permissions`,
                                method: 'POST',
                                data: PERMISSIONS_TO_SEND,
                            }).then((res) => {
                                console.log(res);
                                PERMISSIONS_TO_SEND = [];
                            }).catch((err) => {
                                console.log(err);
                                PERMISSIONS_TO_SEND = [];
                            });
                        }}
                    >
                        OK
                    </Button>
                    <Button
                        marginRight={50}
                        style={{width: '15%'}}
                        appearance='primary'
                        intent='danger'
                        onClick={close}
                    >
                        Back
                    </Button>
                </div>
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
        backgroundColor: 'rgba(0, 0, 0, 0.25)',
        zIndex: 2000,
        alignItems: 'center',
        justifyContent: 'center',
    },
    modal: {
        position: 'relative',
        height: '70%',
        width: '70%',
        padding: '1em',
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
});

export default Permissions;
