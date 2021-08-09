import React from 'react';

export const EditorLoader = () => {
    const [loading, setLoading] = React.useState(true);
    React.useEffect(() => {
        window.addEventListener('ONLYOFFICE_READY', () => setLoading(false));

        return () => window.removeEventListener('ONLYOFFICE_READY', () => setLoading(false));
    }, []);
    return (
        <>
            {loading ? (
                <div
                    className='plugin-editor__loader'
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
