import React from 'react';

import {FileInfo} from 'mattermost-redux/types/files';
import {useDispatch} from 'react-redux';

import {getIconByExt} from 'utils/file';

import {openEditor} from 'redux/actions';

type Props = {
    fileInfo: FileInfo;
};

export default function FilePreviewOverride(props: Props) {
    const dispatch = useDispatch();
    const icon = getIconByExt(props.fileInfo.extension);
    return (
        <div className='modal-image-backround'>
            <div className='modal-image__content'>
                <div className='file-details__container'>
                    <a
                        className='file-details__preview'
                        style={{cursor: 'default'}}
                        onClick={(e) => e.preventDefault()}
                    >
                        <span className='file-details__preview-helper'/>
                        <img
                            alt='file preview'
                            src={icon}
                        />
                    </a>
                    <div
                        className='file-details'
                        style={{position: 'relative'}}
                    >
                        <div className='file-details__name'>{props.fileInfo.name}</div>
                        <div className='file-details__info'>{`File type ${props.fileInfo.extension.toUpperCase()}`}</div>
                        <button
                            className='btn btn-primary'
                            style={{position: 'absolute', right: '2rem', bottom: '2rem'}}
                            type='button'
                            onClick={() => dispatch(openEditor(props.fileInfo))}
                        >
                            {'Open in ONLYOFFICE'}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
