import appDispatcher from '../utils/dispatcher';
import * as http from '../utils/http';
import * as ACTIONS from "./types";
import * as CONSTS from '../utils/consts';

export const getRegisteredUsers = (url) => {
    http.get(url)
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.REGISTERED_USERS_SUCCESS,
            data
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.REGISTERED_USERS_FAIL,
            error
        }))
};

export const getDialerStatus = (url) => {
    http.get(url)
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.DIALER_STATUS_SUCCESS,
            data
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.DIALER_STATUS_FAIL,
            error
        }))
};

export const getCallInProgress = (url) => {
    http.get(url)
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.CALL_IN_PROGRESS_SUCCESS,
            data
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.CALL_IN_PROGRESS_FAIL,
            error
        }))
};

export const startDialer = (url) => {
    http.get(url)
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.DIALER_START_SUCCESS,
            data
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.DIALER_START_FAIL,
            error
        }))
};

export const stopDialer = (url) => {
    http.get(url)
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.DIALER_STOP_SUCCESS,
            data
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.DIALER_STOP_FAIL,
            error
        }))
};

export const pagingChange = (page) => {
    const uri = CONSTS.getHostFn()
        .replace('{API}', CONSTS.CALL_IN_PROGRESS);
    const getParams = CONSTS.LIMIT_PARAMS
        .replace('{limit}', CONSTS.DEFAULT_RECORDS)
        .replace('{page}', page);
    http.get(uri + getParams)
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.PAGING_CHANGE_SUCCESS,
            data,
            page
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.PAGING_CHANGE_FAIL,
            error
        }))
};

export const submitPriority = (id, priority) => {
    const url = CONSTS.getHostFn().replace('{API}', CONSTS.CALL_IN_PROGRESS) + `/${id}`;
    http.put(url, {priority: priority})
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.SUBMIT_CHANGE_PRIORITY_SUCCESS,
            id,
            priority
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.SUBMIT_CHANGE_PRIORITY_FAIL,
            error
        }))
};

export const cancelChangePriority = (id) => {
    appDispatcher.dispatch({
        type: ACTIONS.CHANGE_PRIORITY_CANCEL,
        id
    });
};

export const changePriority = (id) => {
    appDispatcher.dispatch({
        type: ACTIONS.CHANGE_PRIORITY,
        id
    });
};

export const deleteRecord = (id) => {
    const url = CONSTS.getHostFn().replace('{API}', CONSTS.CALL_IN_PROGRESS) + `/${id}`;
    http.del(url)
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.DELETE_RECORD_SUCCESS,
            id
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.DELETE_RECORD_FAIL,
            error
        }))
};

