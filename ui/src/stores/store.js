// utils
import { ReduceStore } from 'flux/utils';
// consts
import * as ACTIONS from '../actions/types';
import appDispatcher from '../utils/dispatcher';
import * as CONSTS from '../utils/consts';


class AppStore extends ReduceStore {

    getInitialState () {
        return {
            operators: [],
            dialerLists: [],
            paging: {
                total: null,
                currentPage: 1,
                numPerPage: CONSTS.DEFAULT_RECORDS,
            },
            isAppStarted: false,
            isAppStopped: false,
            urls: {
                ...CONSTS.API
            },
            sortOrder: CONSTS.ASC,
            sortBy: CONSTS.PRIORITY,
        };
    }

    reduce (state, action) {
        switch (action.type) {

            case ACTIONS.REGISTERED_USERS_SUCCESS:
                return {
                    ...state,
                    operators: action.data.result
                };

            case ACTIONS.CALL_IN_PROGRESS_SUCCESS:
                const callInProgressList = action.data.result.map(item => {
                    item.priority = item.priority.priority;
                    return item;
                });
                return {
                    ...state,
                    paging: {
                        ...state.paging,
                        total: action.data.total,
                    },
                    dialerLists: callInProgressList
                };

            case ACTIONS.DIALER_START_SUCCESS:
                return {
                    ...state,
                    isAppStarted: true,
                    isAppStopped: false,
                };

            case ACTIONS.DIALER_START_FAIL:
                return {
                    ...state,
                    isAppStarted: false,
                    isAppStopped: false,
                };

            case ACTIONS.DIALER_STOP_SUCCESS:
                return {
                    ...state,
                    isAppStarted: false,
                    isAppStopped: true,
                };

            case ACTIONS.DIALER_STOP_FAIL:
                return {
                    ...state,
                    isAppStarted: false,
                    isAppStopped: false,
                };

            case ACTIONS.DIALER_STATUS_SUCCESS:
                const updateStatus = action.data.result ? {
                    isAppStarted: true,
                    isAppStopped: false,
                } : {
                    isAppStarted: false,
                    isAppStopped: true,
                };
                return {
                    ...state,
                    ...updateStatus
                };

            case ACTIONS.SUBMIT_CHANGE_PRIORITY_SUCCESS:
                const updateBysubmitPriority = {
                    dialerLists: [
                        ...state.dialerLists.map(dialer => {
                            if (dialer.id === action.id) {
                                dialer.priority = action.priority;
                            }
                            return dialer;
                        })
                    ]
                };
                return {
                    ...state,
                    ...updateBysubmitPriority
                };

            case ACTIONS.PAGING_CHANGE_SUCCESS:
                const pagingChangeList = action.data.result.map(item => {
                    item.priority = item.priority.priority;
                    return item;
                });
                return {
                    ...state,
                    dialerLists: pagingChangeList,
                    paging: {
                        ...state.paging,
                        currentPage: action.page,
                        total: action.data.total,
                    }
                };

            case ACTIONS.DELETE_RECORD_SUCCESS:
                return {
                    ...state,
                    dialerLists: [
                        ...state.dialerLists.filter(dialer => dialer.id !== action.id)
                    ],
                    paging: {
                        ...state.paging,
                        total: state.paging.total - 1,
                    }
                };

            case ACTIONS.LIMIT_CHANGE_SUCCESS:
                return {
                    ...state,
                    paging: {
                        ...state.paging,
                        numPerPage: action.limit,
                    }
                };

            case ACTIONS.SORT_CHANGE_SUCCESS:
                return {
                    ...state,
                    sortBy: action.sortBy,
                    sortOrder: action.sortOrder === CONSTS.ASC ? CONSTS.DESC : CONSTS.ASC,
                };


            default:
                return state;
        }
    };
}

export default new AppStore(appDispatcher);
