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

import managerDark from 'public/images/manager_dark.svg';
import managerLight from 'public/images/manager_light.svg';
import React from 'react';

type Props = {
    theme: string;
};

export function ManagerIcon({theme}: Props) {
    const icon = theme === 'dark' ? managerDark : managerLight;
    return (
        <img
            style={{width: '16px', height: '16px'}}
            alt={'open manager'}
            src={icon}
        />
    );
}
