import React from 'react';

type UserDetailsProps = {
    username: string
}

export const UserDetails = ({username}: UserDetailsProps) => {
    return (
        <div
            className='more-modal__details'
        >
            <div className='d-flex whitespace--nowrap'>
                <div className='more-modal__name'>
                    <button
                        aria-label={`@${username}`}
                        className='user-popover style--none'
                    >
                        {username}
                    </button>
                </div>
            </div>
        </div>
    );
};
