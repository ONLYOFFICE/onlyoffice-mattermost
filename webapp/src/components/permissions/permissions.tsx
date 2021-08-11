/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable react/jsx-no-literals */
import React, {useState, useEffect} from 'react';
import {Dispatch} from 'redux';
import {FileInfo} from 'mattermost-redux/types/files';

import AsyncSelect from 'react-select/async';
import Select, {OptionTypeBase, OptionsType} from 'react-select';
import makeAnimated from 'react-select/animated';
import {Modal, Button} from 'react-bootstrap';

import {apiGET, apiPOST} from 'api';
import {ONLYOFFICE_PLUGIN_API, ONLYOFFICE_PLUGIN_API_CHANNEL_USER,
    ONLYOFFICE_PLUGIN_API_FILE_PERMISSIONS, ONLYOFFICE_PLUGIN_API_SET_FILE_PERMISSIONS, ONLYOFFICE_WILDCARD_USER} from 'utils';
import {debounce} from 'utils/lodash';
import {FilePermissions, getFileAccess, getPermissionsMap, SubmitPermissionsPayload} from 'utils/file';
import {AutocompleteUser, mapUserToAutocompleteUser, sortAutocompleteUsers, User} from 'utils/user';

import {UserRow} from './permissions_user_row';
import {PermissionsFooter} from './permissions_footer';
import {UserList} from './permissions_user_list';
import {PermissionsHeaderFilter} from './permissions_header_filter';

import 'public/scss/permissions.scss';

type PermissionsProps = {
    visible: boolean,
    close: () => (dispatch: Dispatch) => void,
    fileInfo: FileInfo,
};

const animatedComponents = makeAnimated();

const Permissions: React.FC<PermissionsProps> = ({visible, close, fileInfo}: PermissionsProps) => {
    const [allAccess, setAllAccess] = useState(FilePermissions.READ_ONLY.toString());
    const [current, setCurrent] = useState<AutocompleteUser[]>([]);
    const [users, setUsers] = useState<AutocompleteUser[]>([]);
    const permissionsMap = getPermissionsMap().map((entry: FilePermissions) => {
        return {
            value: entry.toString(),
            label: entry.toString(),
        };
    });

    useEffect(() => {
        if (visible) {
            apiGET(ONLYOFFICE_PLUGIN_API + ONLYOFFICE_PLUGIN_API_FILE_PERMISSIONS + fileInfo.id).then((resUser: User[]) => {
                if (!resUser) {
                    return;
                }
                const permissions: AutocompleteUser[] = [];
                // eslint-disable-next-line max-nested-callbacks
                resUser.forEach((user: User) => {
                    const mappedUser = mapUserToAutocompleteUser(user);
                    if (user.id === ONLYOFFICE_WILDCARD_USER) {
                        setAllAccess(mappedUser.permissions);
                    } else {
                        permissions.push(mappedUser);
                    }
                });
                sortAutocompleteUsers(permissions);
                setUsers(permissions);
            }).catch();
        }
    }, [visible, fileInfo]);

    if (!visible) {
        return null;
    }

    const load = debounce((input: any, callback: any) => {
        if (!input) {
            return;
        }
        if (users.find((user: AutocompleteUser) => user.label === input)) {
            callback([]);
            return;
        }
        apiGET(ONLYOFFICE_PLUGIN_API + ONLYOFFICE_PLUGIN_API_CHANNEL_USER + input, {
            ONLYOFFICE_FILEID: fileInfo.id,
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
    }, 2000);

    const onChange = (value: OptionTypeBase | OptionsType<OptionTypeBase> | null) => {
        setCurrent((value as AutocompleteUser[]));
    };

    const onAllChange = (value: any) => {
        if (value.label) {
            setAllAccess(value.label);
        }
    };

    const onExit = () => {
        setAllAccess(FilePermissions.READ_ONLY.toString());
        setCurrent([]);
        setUsers([]);

        const modal = document.getElementById('onlyoffice-permissions-modal');
        const backdrop = modal?.previousElementSibling;

        // eslint-disable-next-line no-unused-expressions
        modal?.classList.remove('in');
        // eslint-disable-next-line no-unused-expressions
        backdrop?.classList.remove('in');

        setTimeout(() => close(), 300);
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

    const onSubmitPermissions = () => {
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

        apiPOST(ONLYOFFICE_PLUGIN_API + ONLYOFFICE_PLUGIN_API_SET_FILE_PERMISSIONS, JSON.stringify(payload)).then(() => {
            onExit();
        }).catch();
    };

    return (
        <Modal
            show={visible}
            onHide={onExit}
            onExited={onExit}
            role='dialog'
            id='onlyoffice-permissions-modal'
        >
            <Modal.Header
                closeButton={true}
            >
                {`Sharing Settings ${fileInfo.name.split('.')[0]}`}
                <button
                    type='button'
                    className='close'
                    aria-label='Close'
                    onClick={onExit}
                >
                    <span aria-hidden='true'>×</span>
                    <span className='sr-only'>Close</span>
                </button>
            </Modal.Header>
            <div className='onlyoffice-permissions-modal__body'>
                <div className='filtered-user-list'>
                    <div
                        className='filter-row'
                        style={{marginBottom: '1rem', marginTop: '1rem'}}
                    >
                        <PermissionsHeaderFilter>
                            <div style={{flexGrow: 1, marginRight: '2rem'}}>
                                <AsyncSelect
                                    id='onlyoffice-permissions-select'
                                    className='react-select-container'
                                    classNamePrefix='react-select'
                                    closeMenuOnSelect={false}
                                    components={animatedComponents}
                                    isMulti={true}
                                    loadOptions={load}
                                    onChange={onChange}
                                    value={current}
                                />
                            </div>
                            <Button
                                className='btn btn-md btn-primary'
                                disabled={current.length === 0}
                                onClick={() => {
                                    if (current) {
                                        const contentSection = document.getElementById('scroller-dummy');
                                        setTimeout(() => contentSection?.scrollIntoView({behavior: 'smooth'}), 300);
                                        setUsers((prevUsers: AutocompleteUser[]) => [...prevUsers, ...current]);
                                        setCurrent([]);
                                    }
                                }}
                            >
                                Add
                            </Button>
                        </PermissionsHeaderFilter>
                        <div
                            className='col-sm-12'
                            style={{marginTop: '2rem', display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}
                        >
                            <span
                                style={{flexGrow: 2}}
                                className='member-count pull-left'
                            >
                                <span>Default access rights for chat members</span>
                            </span>
                            <div style={{marginRight: '2.5rem', width: '10rem'}}>
                                <Select
                                    isSearchable={false}
                                    value={{
                                        value: allAccess,
                                        label: allAccess,
                                    }}
                                    options={permissionsMap}
                                    onChange={onAllChange}
                                />
                            </div>
                        </div>
                    </div>
                    <UserList>
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
                        <div id='scroller-dummy'/>
                    </UserList>
                    <PermissionsFooter>
                        <Button
                            className='btn btn-md'
                            style={{marginRight: '1rem', border: 'none'}}
                            onClick={onExit}
                        >
                            <span style={{color: '#2389D7'}}>Cancel</span>
                        </Button>
                        <Button
                            className='btn btn-md btn-primary'
                            onClick={onSubmitPermissions}
                        >
                            Save
                        </Button>
                    </PermissionsFooter>
                </div>
            </div>
        </Modal>
    );
};

export default Permissions;
