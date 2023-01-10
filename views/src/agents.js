import superagent from 'superagent'

const API_ROOT = "/api"

const responseBody = (res) => ({
    res: res.body,
    raw: res,
});

const error = (err) => ({
    res: err.response.body,
    raw: err.response,
});

const tokenPlugin = (req) => {};

const requests = {
    del: (url) =>
        superagent.del(`${API_ROOT}${url}`).use(tokenPlugin).then(responseBody),
    get: (url) =>
        superagent.get(`${API_ROOT}${url}`).use(tokenPlugin).then(responseBody),
    put: (url, body) =>
        superagent
            .put(`${API_ROOT}${url}`, body)
            .use(tokenPlugin)
            .then(responseBody),
    post: (url, body) =>
        superagent
            .post(`${API_ROOT}${url}`, body)
            .use(tokenPlugin)
            .then(responseBody)
            .catch(error),
};


export default {
    userprofiles: () => requests.get('/userprofiles'),
    servers: () => requests.get('/servers'),
}
