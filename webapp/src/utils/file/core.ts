import {FileInfo} from 'mattermost-redux/types/files';

import xls from 'public/images/xls.svg';
import xlsx from 'public/images/xlsx.svg';
import csv from 'public/images/csv.svg';
import doc from 'public/images/doc.svg';
import docx from 'public/images/docx.svg';
import ppt from 'public/images/ppt.svg';
import pptx from 'public/images/pptx.svg';
import word from 'public/images/generic_word.svg';
import cell from 'public/images/generic_cell.svg';
import slide from 'public/images/generic_slide.svg';

const ONLYOFFICE_CELL = 'cell';
const ONLYOFFICE_WORD = 'word';
const ONLYOFFICE_SLIDE = 'slide';

const AllowedExtensionsMap = new Map([
    ['xls', ONLYOFFICE_CELL],
    ['xlsx', ONLYOFFICE_CELL],
    ['xlsm', ONLYOFFICE_CELL],
    ['xlt', ONLYOFFICE_CELL],
    ['xltx', ONLYOFFICE_CELL],
    ['xltm', ONLYOFFICE_CELL],
    ['ods', ONLYOFFICE_CELL],
    ['fods', ONLYOFFICE_CELL],
    ['ots', ONLYOFFICE_CELL],
    ['csv', ONLYOFFICE_CELL],
    ['pps', ONLYOFFICE_SLIDE],
    ['ppsx', ONLYOFFICE_SLIDE],
    ['ppsm', ONLYOFFICE_SLIDE],
    ['ppt', ONLYOFFICE_SLIDE],
    ['pptx', ONLYOFFICE_SLIDE],
    ['pptm', ONLYOFFICE_SLIDE],
    ['pot', ONLYOFFICE_SLIDE],
    ['potx', ONLYOFFICE_SLIDE],
    ['potm', ONLYOFFICE_SLIDE],
    ['odp', ONLYOFFICE_SLIDE],
    ['fodp', ONLYOFFICE_SLIDE],
    ['otp', ONLYOFFICE_SLIDE],
    ['doc', ONLYOFFICE_WORD],
    ['docx', ONLYOFFICE_WORD],
    ['docm', ONLYOFFICE_WORD],
    ['dot', ONLYOFFICE_WORD],
    ['dotx', ONLYOFFICE_WORD],
    ['dotm', ONLYOFFICE_WORD],
    ['odt', ONLYOFFICE_WORD],
    ['fodt', ONLYOFFICE_WORD],
    ['ott', ONLYOFFICE_WORD],
    ['rtf', ONLYOFFICE_WORD],
]);

const ExtensionIcons = new Map([
    ['xls', xls],
    ['xlsx', xlsx],
    ['csv', csv],
    ['ppt', ppt],
    ['pptx', pptx],
    ['doc', doc],
    ['docx', docx],
    [ONLYOFFICE_WORD, word],
    [ONLYOFFICE_CELL, cell],
    [ONLYOFFICE_SLIDE, slide],
]);

export function getIconByExt(fileExt: string): string {
    if (ExtensionIcons.has(fileExt)) {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        return ExtensionIcons.get(fileExt)!;
    }
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return ExtensionIcons.get(getFileTypeByExt(fileExt))!;
}

export function getFileTypeByExt(fileExt: string): string {
    if (AllowedExtensionsMap.has(fileExt)) {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        return AllowedExtensionsMap.get(fileExt)!;
    }
    return '';
}

export function isExtensionSupported(fileExt: string): boolean {
    if (AllowedExtensionsMap.has(fileExt)) {
        return true;
    }

    return false;
}

export function isFileAuthor(fileInfo: FileInfo): boolean {
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
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        return parts!.pop()!.split(';').shift() || '';
    }
    return '';
}

