// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/**
 *
 * (c) Copyright Ascensio System SIA 2025
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

import { debounceUsersLoad } from 'util/func';
import { getTranslations } from 'util/lang';
import type { FileAccess } from 'util/permission';
import { getFileAccess } from 'util/permission';
import type { MattermostUser } from 'util/user';

import React, { useState, useEffect } from 'react';
import { Button } from 'react-bootstrap';
import Select from 'react-select';
import AsyncSelect from 'react-select/async';

import type { Channel } from 'mattermost-redux/types/channels';
import type { FileInfo } from 'mattermost-redux/types/files';

type Props = {
  loading: boolean;
  channel: Channel | undefined | null;
  fileInfo: FileInfo;
  wildcardAccess: string;
  users: MattermostUser[];
  onSetWildcardAccess: (value: any) => void;
  onAppendUsers: (newUsers: MattermostUser[]) => void;
};

export const PermissionsHeader = (props: Props) => {
  const i18n = getTranslations();
  const permissionsMap = getFileAccess().map((entry: FileAccess) => {
    return {
      value: entry.toString(),
      label: i18n[`types.permissions.${entry.toString().toLowerCase() as 'edit' | 'read'}`] || entry.toString(),
    };
  });
  const [current, setCurrent] = useState<MattermostUser[]>([]);
  const [accessHeader, setAccessHeader] = useState<string>(i18n['permissions.loading']);

  useEffect(() => {
    const isChannel = window.location.href.split('/').includes('channels');
    if (!props.loading) {
      setAccessHeader(isChannel ? i18n['permissions.access_header_default'] : i18n['permissions.access_header']);
    }

    return () => setCurrent([]);
  }, [props.channel, props.loading]);

  return (
    <div
      className='filter-row'
      style={props.channel ? { marginBottom: '1rem', marginTop: '1rem' } : { maxHeight: '10rem' }}
    >
      {props.channel && (
        <div
          className='col-xs-12'
          style={{ marginBottom: '1rem' }}
        >
          <div style={{ display: 'flex' }}>
            <div style={{ flexGrow: 1, marginRight: '0.5rem' }}>
              <AsyncSelect
                id='onlyoffice-permissions-select'
                placeholder={i18n['permissions.modal_search_placeholder']}
                loadingMessage={() => i18n['permissions.modal_search_loading_placeholder']}
                noOptionsMessage={() => i18n['permissions.modal_search_no_options_placeholder']}
                className='react-select-container'
                classNamePrefix='react-select'
                closeMenuOnSelect={false}
                isMulti={true}
                loadOptions={debounceUsersLoad(props.channel, props.fileInfo, props.users)}
                onChange={(users) => setCurrent((users as MattermostUser[]))}
                value={current}
                isDisabled={props.loading || !props.channel}
                components={{
                  DropdownIndicator: () => null,
                  IndicatorSeparator: () => null,
                }}
                styles={{
                  container: (provided: any) => ({ ...provided, height: '100%' }),
                  control: (provided: any) => ({ ...provided, minHeight: '100%' }),
                  multiValue: (provided: any) => ({
                    ...provided,
                    backgroundColor: '#f0f0f0',
                    borderRadius: '49px',
                    margin: '2px 4px',
                    padding: '2px 4px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '6px',
                  }),
                  multiValueLabel: (provided: any) => ({
                    ...provided,
                    textAlign: 'center',
                    color: '#3d3c40',
                    fontSize: '12px',
                    fontWeight: 400,
                    lineHeight: '16px',
                    padding: 0,
                  }),
                  multiValueRemove: (provided: any) => ({
                    ...provided,
                    width: '10px',
                    height: '10px',
                    minWidth: '10px',
                    minHeight: '10px',
                    borderRadius: '50%',
                    margin: 0,
                    padding: 0,
                    fontSize: '0.8rem',
                    lineHeight: 1,
                    border: '1px solid #ababad',
                    backgroundColor: '#ababad',
                    color: '#f0f0f0',
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    ':hover': {
                      backgroundColor: '#9c9c9e',
                    },
                  }),
                }}
              />
            </div>
            <Button
              className='btn btn-md btn-primary'
              disabled={current.length === 0 || props.loading}
              onClick={() => {
                if (current) {
                  const contentSection = document.getElementById('scroller-dummy');
                  setTimeout(() => contentSection?.scrollIntoView({ behavior: 'smooth' }), 300);
                  props.onAppendUsers(current);
                  setCurrent([]);
                }
              }}
            >
              {i18n['permissions.modal_button_add']}
            </Button>
          </div>
        </div>
      )}
      <div
        className='col-sm-12'
        style={{ marginTop: '2rem', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}
      >
        <span
          className='member-count pull-left onlyoffice-permissions__access-header'
        >
          <span>{accessHeader}</span>
        </span>
        <div style={{ marginLeft: '10px' }}>
          <Select
            isSearchable={false}
            value={{
              value: props.wildcardAccess,
              label: i18n[`types.permissions.${props.wildcardAccess.toLowerCase() as 'edit' | 'read'}`] || props.wildcardAccess,
            }}
            options={permissionsMap}
            onChange={(selected) => props.onSetWildcardAccess(selected?.value)}
            isDisabled={props.loading}
            components={{
              IndicatorSeparator: () => null,
            }}
            styles={{
              control: (provided: any) => ({
                ...provided,
                width: 'auto',
                height: '32px',
                border: 'none',
                borderRadius: '4px',
                boxShadow: 'none',
                backgroundColor: 'transparent',
                cursor: 'pointer',
                display: 'flex',
                justifyContent: 'flex-end',
                padding: '4px 10px 5px 12px',
                ':hover': {
                  backgroundColor: '#1C58D914',
                },
              }),
              valueContainer: (provided: any) => ({
                ...provided,
                padding: 0,
                display: 'flex',
                justifyContent: 'flex-end',
              }),
              indicatorsContainer: (provided: any) => ({
                ...provided,
                padding: 0,
                display: 'flex',
                justifyContent: 'flex-end',
              }),
              singleValue: (provided: any) => ({
                ...provided,
                color: '#1C58D9',
                marginRight: '6px',
              }),
              dropdownIndicator: (provided: any) => ({
                ...provided,
                color: '#1C58D9',
                padding: 0,
                marginRight: '0px',
                ':hover': {
                  color: '#1C58D9',
                },
                svg: {
                  width: '14px',
                  height: '14px',
                  fill: '#1C58D9',
                  ':hover': {
                    fill: '#1C58D9',
                  },
                }
              }),
            }}
          />
        </div>
      </div>
    </div>
  );
};
