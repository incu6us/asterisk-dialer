import appDispatcher from '../utils/dispatcher';
import * as http from '../utils/http';
import * as ACTIONS from "./types";
import * as consts from '../utils/consts';

export const getRegisteredUsers = (url) =
>
{
    http.get(url)
        .then(data = > appDispatcher.dispatch({
        type: ACTIONS.REGISTERED_USERS_SUCCESS,
        data
    })
)
.
    catch(error = > appDispatcher.dispatch({
        type: ACTIONS.REGISTERED_USERS_FAIL,
        error
    })
)
}
;

export const getDialerStatus = (url) =
>
{
    http.get(url)
        .then(data = > appDispatcher.dispatch({
        type: ACTIONS.DIALER_STATUS_SUCCESS,
        data
    })
)
.
    catch(error = > appDispatcher.dispatch({
        type: ACTIONS.DIALER_STATUS_FAIL,
        error
    })
)
}
;

export const startDialer = (url) =
>
{
    http.get(url)
        .then(data = > appDispatcher.dispatch({
        type: ACTIONS.DIALER_START_SUCCESS,
        data
    })
)
.
    catch(error = > appDispatcher.dispatch({
        type: ACTIONS.DIALER_START_FAIL,
        error
    })
)
}
;

export const stopDialer = (url) =
>
{
    http.get(url)
        .then(data = > appDispatcher.dispatch({
        type: ACTIONS.DIALER_STOP_SUCCESS,
        data
    })
)
.
    catch(error = > appDispatcher.dispatch({
        type: ACTIONS.DIALER_STOP_FAIL,
        error
    })
)
}
;

export const pagingChange = (page, url) =
>
{
    http.get(url)
        .then(data = > appDispatcher.dispatch({
        type: ACTIONS.PAGING_CHANGE_SUCCESS,
        data,
        page
    })
)
.
    catch(error = > appDispatcher.dispatch({
        type: ACTIONS.PAGING_CHANGE_FAIL,
        error
    })
)
}
;

