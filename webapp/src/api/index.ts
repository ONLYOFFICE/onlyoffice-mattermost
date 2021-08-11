export const apiGET = async (url: string, headers?: HeadersInit) => {
    let json;
    const response = await fetch(url, {
        method: 'GET',
        headers,
    });

    if (response.body) {
        json = await response.json();
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
