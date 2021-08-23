import React from 'react';

type PermissionsFooterProps = {
    children: React.ReactNode
}

export const PermissionsFooter = ({children}: PermissionsFooterProps) => {
    return (
        <div
            className='filter-controls'
            style={{display: 'flex', justifyContent: 'flex-end', padding: 0, margin: '1rem'}}
        >
            {children}
        </div>
    );
};
