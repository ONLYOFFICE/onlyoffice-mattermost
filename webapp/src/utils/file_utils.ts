import {FileInfo} from 'mattermost-redux/types/files';

const AllowedExtensions = [
    'xls',
    'xlsx',
    'xlsm',
    'xlt',
    'xltx',
    'xltm',
    'ods',
    'fods',
    'ots',
    'csv',
    'pps',
    'ppsx',
    'ppsm',
    'ppt',
    'pptx',
    'pptm',
    'pot',
    'potx',
    'potm',
    'odp',
    'fodp',
    'otp',
    'doc',
    'docx',
    'docm',
    'dot',
    'dotx',
    'dotm',
    'odt',
    'fodt',
    'ott',
    'rtf',
    'txt',
    'html',
    'htm',
    'mht',
    'pdf',
    'djvu',
    'fb2',
    'epub',
    'xps',
];

export function isExtensionSupported(fileExt: string): boolean {
    if (AllowedExtensions.find((ext) => ext === fileExt)) {
        return true;
    }

    return false;
}

export function isFileAuthor(fileInfo: FileInfo): boolean {
    // eslint-disable-next-line no-console
    const userId: string = getCookie('MMUSERID');

    if (userId) {
        return fileInfo.user_id === userId;
    }

    return false;
}

function getCookie(name: string): string {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        return parts!.pop()!.split(';').shift() || '';
    }
    return '';
}

