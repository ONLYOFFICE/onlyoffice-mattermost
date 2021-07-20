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

