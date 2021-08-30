/**
 *
 * (c) Copyright Ascensio System SIA 2021
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

import React from 'react';

type UserItemProps = {
    alt: string,
    src: string,
}

export const UserIcon = (props: UserItemProps) => {
    return (
        <button
            className='statuc-wrapper style--none'
            tabIndex={-1}
        >
            <span className='profile-icon'>
                <img
                    className='Avatar Avatar-md'
                    alt={props.alt}
                    src={props.src}
                />
            </span>
        </button>
    );
};
