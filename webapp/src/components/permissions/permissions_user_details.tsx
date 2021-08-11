import React from 'react';

type UserDetailsProps = {
    username: string
    email: string
}

export const UserDetails = ({username, email}: UserDetailsProps) => {
    return (
        <div
            className='more-modal__details'
        >
            <div className='more-modal__name'>
                <button
                    aria-label={`@${username}`}
                    className='user-popover style--none'
                >
                    {`@${username}`}
                </button>
            </div>
            <div className='more-modal__description'>
                {email}
            </div>
        </div>
    );
};
