import React from 'react';

type PermissionsHeaderFilterProps = {
    children: React.ReactNode
}

export const PermissionsHeaderFilter = ({children}: PermissionsHeaderFilterProps) => {
    return (
        <div
            className='col-xs-12'
            style={{marginBottom: '1rem'}}
        >
            <div style={{display: 'flex'}}>
                {children}
            </div>
        </div>
    );
};
