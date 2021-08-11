import React from 'react';

type UserActionsProps = {
    children: React.ReactNode
}

export const UserActions = ({children}: UserActionsProps) => {
    return (
        <div
            className='more-modal__actions'
            style={{display: 'flex', paddingRight: '0.3rem'}}
        >
            {children}
        </div>
    );
};
