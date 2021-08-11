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
