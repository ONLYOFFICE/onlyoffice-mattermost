import { useState, useEffect } from 'react';
import axios from 'axios';

export function useSearchUsers(pageNumber, limit, fileInfoId) {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);
    const [hasMore, setHasMore] = useState(false);

    useEffect(() => {
        setLoading(true);
        setError(false);
        let cancel;
        axios({
            url: `http://46.101.101.37:8065/plugins/com.onlyoffice.mattermost-plugin/onlyofficeapi/channel_users`,
            method: 'GET',
            headers: {
                ONLYOFFICE_FILEID: fileInfoId,
            },
            params: {
                page: pageNumber,
                limit
            },
            cancelToken: new axios.CancelToken(c => cancel = c)
        }).then((res) => {
            if (res.data) {
                setData(prevData => {
                    return [...new Set([...prevData, ...res.data])];
                });
                setHasMore(res.data.length > 0);
                setLoading(false);
            } else {
                setHasMore(false);
                setLoading(false);
            }
        }).catch((err) => {
            if (axios.isCancel(err)) return;
            setError(true);
        });
        return () => cancel();
    }, [pageNumber, limit, fileInfoId]);

    return { loading, error, data, hasMore };
}
