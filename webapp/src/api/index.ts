export const apiGET = async <T>(url: string, headers?: HeadersInit): Promise<T> => {
    let json: T;

    try {
        const response = await fetch(url, {
            method: 'GET',
            headers,
        });

        json = await response.json();
    } catch {
        throw new Error('API GET call error');
    }

    return json;
};

export const apiPOST = async (url: string, body: string, headers?: HeadersInit) => {
    try {
        await fetch(url, {
            method: 'POST',
            headers,
            body,
        });
    } catch {
        throw new Error('API POST call error');
    }
};
