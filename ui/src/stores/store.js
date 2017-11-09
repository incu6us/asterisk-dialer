// utils
import { ReduceStore } from 'flux/utils';
// consts
import * as ACTIONS from '../actions/types';
import appDispatcher from '../utils/dispatcher';


class AppStore extends ReduceStore {

    getInitialState () {
        return {
            operators: [],
            dialerLists: [],
            paging: {
                total: null,
                currentPage: 1,
                numPerPage: 10,
            },
            isAppStarted: false,
            isAppStopped: false,
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
                console.log('action ->>', action.data.result);
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

            case ACTIONS.CHANGE_PRIORITY:
                const updateByChangePriority = {
                    dialerLists: [
                        ...state.dialerLists.map(dialer => {
                            if (dialer.msisdn === action.msisdn) {
                                dialer.isChanging = true;
                            }
                            return dialer;
                        })
                    ]
                };
                return {
                    ...state,
                    ...updateByChangePriority
                };

            case ACTIONS.CHANGE_PRIORITY_SUBMIT:
                const updateBysubmitPriority = {
                    dialerLists: [
                        ...state.dialerLists.map(dialer => {
                            if (dialer.msisdn === action.msisdn) {
                                dialer.isChanging = false;
                            }
                            return dialer;
                        })
                    ]
                };
                return {
                    ...state,
                    ...updateBysubmitPriority
                };

            case ACTIONS.CHANGE_PRIORITY_CANCEL:
                const updateByCancelPriority = {
                    dialerLists: [
                        ...state.dialerLists.map(dialer => {
                            if (dialer.msisdn === action.msisdn) {
                                dialer.isChanging = false;
                            }
                            return dialer;
                        })
                    ]
                };
                return {
                    ...state,
                    ...updateByCancelPriority
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

            default:
                return state;
        };
    };
}

export default new AppStore(appDispatcher);
