export function getUrl (baseUrl, page, limit=20, sortBy, sortOrder) {
    const params = {
        page: page,
        limit: limit,
        sortBy: sortBy,
        sortOrder: sortOrder,
    };
    return baseUrl + '?' + Object.keys(params).map(param => {
        const value = params[param];
        return value ? param + '=' + encodeURIComponent(value): null;
    }).filter(item => item !== null).join('&');
}
