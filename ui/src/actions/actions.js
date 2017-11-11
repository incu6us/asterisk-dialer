import appDispatcher from '../utils/dispatcher';
import * as http from '../utils/http';
import * as ACTIONS from "./types";

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

export const pagingChange = (url, page) => {
    http.get(url)
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

export const submitPriority = (url, id, priority) => {
    http.put(url + `/${id}`, {priority: priority})
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

export const deleteRecord = (url, id) => {
    http.del(url + `/${id}`)
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.DELETE_RECORD_SUCCESS,
            id
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.DELETE_RECORD_FAIL,
            error
        }))
};

export const limitChange = (url, limit) => {
    http.get(url)
        .then(data => appDispatcher.dispatch({
            type: ACTIONS.LIMIT_CHANGE_SUCCESS,
            limit
        }))
        .catch(error => appDispatcher.dispatch({
            type: ACTIONS.LIMIT_CHANGE_FAIL,
            error
        }))
};

