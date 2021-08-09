/* eslint-disable react/jsx-no-literals */
import React, {useState, useEffect} from 'react';

import {FileInfo} from 'mattermost-redux/types/files';

import {Dispatch} from 'redux';

import 'public/scss/classes/onlyoffice_permissions.scss';
import {Modal, Button} from 'react-bootstrap';
import AsyncSelect from 'react-select/async';
import makeAnimated from 'react-select/animated';

import {OptionTypeBase, OptionsType, ActionMeta} from 'react-select';

import {FileAccess, FilePermissions, getFileAccess, getPermissionsMap, SubmitPermissionsPayload} from 'utils/file';

import {debounce} from 'utils/lodash';

import {id as pluginName} from 'manifest';
import {AutocompleteUser, mapUserToAutocompleteUser, User} from 'utils/user';

import {UserRow} from './permissions_user';

const animatedComponents = makeAnimated();

type PermissionsProps = {
    visible: boolean,
    close: () => (dispatch: Dispatch) => void,
    fileInfo: FileInfo,
};

const Permissions: React.FC<PermissionsProps> = ({visible, close, fileInfo}: PermissionsProps) => {
    const [allAccess, setAllAccess] = useState(FilePermissions.READ_ONLY.toString());
    const [current, setCurrent] = useState<AutocompleteUser[]>([]);
    const [users, setUsers] = useState<AutocompleteUser[]>([]);
    const permissionsMap = getPermissionsMap();

    useEffect(() => {
        if (visible) {
            fetch(`/plugins/${pluginName}/onlyofficeapi/get_file_permissions?fileId=${fileInfo.id}`, {
                method: 'GET',
            }).then((res) => {
                return res.json();
            }).then((resUser: User[]) => {
                const permissions: AutocompleteUser[] = [];
                // eslint-disable-next-line max-nested-callbacks
                resUser.forEach((user: User) => {
                    const mappedUser = mapUserToAutocompleteUser(user);
                    if (user.id === '*') {
                        setAllAccess(mappedUser.permissions);
                    } else {
                        permissions.push(mappedUser);
                    }
                });
                setUsers(permissions);
            }).catch((err) => {
                console.log(err);
            });
        }
    }, [visible]);

    if (!visible) {
        return null;
    }

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const load = debounce((input: any, callback: any) => {
        if (!input) {
            return;
        }
        fetch(`/plugins/${pluginName}/onlyofficeapi/channel_user?username=${input}`, {
            method: 'GET',
            headers: {
                ONLYOFFICE_FILEID: fileInfo.id,
            },
        }).then((res) => {
            return res.json();
        }).then((resUser: User) => {
            if (!resUser.id) {
                callback([]);
                return;
            }
            const user = mapUserToAutocompleteUser(resUser);
            callback([user]);
        }).catch(() => {
            callback([]);
        });
    }, 3000);

    const onChange = (value: OptionTypeBase | OptionsType<OptionTypeBase> | null, action: ActionMeta<OptionTypeBase>) => {
        setCurrent((value as AutocompleteUser[]));
    };

    const handleOnExit = () => {
        setAllAccess(FilePermissions.READ_ONLY.toString());
        setCurrent([]);
        setUsers([]);
        close();
    };

    const onRemoveUser = (username: string) => {
        setUsers((prevUsers: AutocompleteUser[]) => prevUsers.filter((user: AutocompleteUser) => user.label !== username));
    };

    const onChangeUserPermissions = (username: string, newPermissions: string) => {
        setUsers((prevUsers: AutocompleteUser[]) => prevUsers.map((user: AutocompleteUser) => {
            if (user.label === username) {
                user.permissions = newPermissions;
            }
            return user;
        }));
    };

    const handleSubmitPermissions = () => {
        const payload: SubmitPermissionsPayload[] = [];
        const allUsers: SubmitPermissionsPayload = {
            FileId: fileInfo.id,
            Username: '*',
            Permissions: FilePermissions.READ_ONLY.toString() === allAccess ? getFileAccess(FilePermissions.READ_ONLY) : getFileAccess(FilePermissions.EDIT_ONLY),
        };
        payload.push(allUsers);
        users.forEach((user: AutocompleteUser) => {
            payload.push({
                FileId: fileInfo.id,
                Username: user.label,
                Permissions: FilePermissions.READ_ONLY.toString() === user.permissions ? getFileAccess(FilePermissions.READ_ONLY) : getFileAccess(FilePermissions.EDIT_ONLY),
            });
        });

        fetch(`/plugins/${pluginName}/onlyofficeapi/set_file_permissions`, {
            method: 'POST',
            body: JSON.stringify(payload),
        }).then((res) => {
            console.log(res);
            close();
        }).catch((err) => {
            console.log(err);
        });
    };

    return (
        <Modal
            show={visible}
            onHide={handleOnExit}
            onExited={handleOnExit}
            role='dialog'
            className='permissions-modal'
        >
            <Modal.Header
                closeButton={true}
            >
                {`Sharing Settings ${fileInfo.name.split('.')[0]}`}
                <button
                    type='button'
                    className='close'
                    aria-label='Close'
                    onClick={handleOnExit}
                >
                    <span aria-hidden='true'>×</span>
                    <span className='sr-only'>Close</span>
                </button>
            </Modal.Header>
            <div>
                <div className='filtered-user-list'>
                    <div
                        className='filter-row'
                        style={{marginBottom: '1rem', marginTop: '1rem'}}
                    >
                        <div
                            className='col-xs-12'
                            style={{marginBottom: '1rem'}}
                        >
                            <div style={{display: 'flex'}}>
                                <div style={{flexGrow: 1, marginRight: '2rem'}}>
                                    <AsyncSelect
                                        closeMenuOnSelect={false}
                                        components={animatedComponents}
                                        isMulti={true}
                                        loadOptions={load}
                                        onChange={onChange}
                                        value={current}
                                    />
                                </div>
                                <Button
                                    style={{backgroundColor: '#166DE0', color: '#FFFFFF', border: 'none'}}
                                    disabled={current.length === 0}
                                    onClick={() => {
                                        if (current) {
                                            setUsers((prevUsers: AutocompleteUser[]) => [...new Set([...prevUsers, ...current])]);
                                            setCurrent([]);
                                        }
                                    }}
                                >Invite</Button>
                            </div>
                        </div>
                        <div
                            className='col-sm-12'
                            style={{marginTop: '2rem', display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}
                        >
                            <span
                                className='member-count pull-left'
                            >
                                <span>Default access rights for chat members</span>
                            </span>
                            <select
                                style={{marginRight: '2.5rem'}}
                                value={allAccess}
                                onChange={(e) => setAllAccess(e.target.value)}
                            >
                                {permissionsMap.map((value: [FilePermissions, FileAccess]) => {
                                    const permString = value[0].toString();
                                    return (
                                        <option
                                            key={permString}
                                            value={permString}
                                        >{permString}</option>
                                    );
                                })}
                            </select>
                        </div>
                    </div>
                    <div className='more-modal__list'>
                        <div>
                            {users.map((user: AutocompleteUser) => {
                                return (
                                    <UserRow
                                        key={user.value}
                                        user={user}
                                        changePermissions={onChangeUserPermissions}
                                        removeUser={onRemoveUser}
                                    />
                                );
                            })}
                        </div>
                    </div>
                    <div
                        className='filter-controls'
                        style={{display: 'flex', justifyContent: 'flex-end', padding: 0, margin: '1rem'}}
                    >
                        <Button
                            style={{backgroundColor: 'white', border: 'none'}}
                            onClick={() => close()}
                        >
                            <span style={{color: '#2389D7'}}>Cancel</span>
                        </Button>
                        <Button
                            style={{backgroundColor: '#166DE0', color: '#FFFFFF', border: 'none'}}
                            onClick={() => {
                                handleSubmitPermissions();
                            }}
                        >
                            Save
                        </Button>
                    </div>
                </div>
            </div>
        </Modal>
    );
};

export default Permissions;
