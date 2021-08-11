import React from 'react';

export const EditorLoader = () => {
    const [loading, setLoading] = React.useState(true);

    const disableLoading = () => {
        setLoading(false);
    };

    React.useEffect(() => {
        window.addEventListener('ONLYOFFICE_READY', disableLoading);

        return () => window.removeEventListener('ONLYOFFICE_READY', disableLoading);
    }, []);
    return (
        <>
            {loading ? (
                <div
                    className='onlyoffice-editor__loader'
                    id='editor-spinner'
                >
                    <></>
                </div>
            ) : (
                null
            )}
        </>
    );
};
