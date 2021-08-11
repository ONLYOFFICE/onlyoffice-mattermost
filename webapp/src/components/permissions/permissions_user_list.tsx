import React from 'react';

type UserListProps = {
    children: React.ReactNode
}

export const UserList = ({children}: UserListProps) => {
    return (
        <div className='more-modal__list'>
            <div>
                {children}
            </div>
        </div>
    );
};
