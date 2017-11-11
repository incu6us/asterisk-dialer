export function getUrl (baseUrl, page, limit=20) {
    const params = {
        page: page,
        //sortBy: sortBy,
        //sortOrder: sortOrder,
        limit: limit,
    };
    return baseUrl + '?' + Object.keys(params).map(param => {
        const value = params[param];
        return value ? param + '=' + encodeURIComponent(value): null;
    }).filter(item => item !== null).join('&');
}
