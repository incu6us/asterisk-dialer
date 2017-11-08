// utils
import { ReduceStore } from 'flux/utils';
// consts
import * as ACTION from '../actions/types';
import appDispatcher from '../utils/dispatcher';


class AppStore extends ReduceStore {

    getInitialState () {
        return {
            operators: [
                {
                    action: "Hangup",
                    exten: "0937530213",
                    peer: "1001",
                    peerStatus: "Registered",
                    time: "2017-11-07T22:24:56+02:00"
                },
                {
                    action: "Fuckup",
                    exten: "0937530213",
                    peer: "1002",
                    peerStatus: "Registered",
                    time: "2017-11-07T22:25:56+02:00"
                },
                {
                    action: "Lookup",
                    exten: "0937530213",
                    peer: "1002",
                    peerStatus: "Fucking shit",
                    time: "2017-11-07T22:25:56+02:00"
                }
            ],
            dialerLists: [
                {
                    msisdn: '0633467991',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1017',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '6',
                },
                {
                    msisdn: '0633467992',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1018',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '7',
                },
                {
                    msisdn: '0633467993',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1019',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '7'
                },
                {
                    msisdn: '0633467994',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1019',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '7'
                },
                {
                    msisdn: '0633467995',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1019',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '7'
                },
                {
                    msisdn: '0633467996',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1019',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '7'
                },
                {
                    msisdn: '0633467997',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1019',
                    timeCalled: '2017-11-07 09:12:08'
                },
                {
                    msisdn: '0633467998',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1019',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '7'
                },
                {
                    msisdn: '0633467999',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1019',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '7'
                },
                {
                    msisdn: '0633467987',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1019',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '7'
                },
                {
                    msisdn: '0633467988',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1017',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '6',
                },
                {
                    msisdn: '0633467989',
                    status: 'progress',
                    time: '2017-11-07 09:00:01',
                    causeTxt: 'User alerting, no answer',
                    event: 'Hangup',
                    callerIdNum: '1018',
                    timeCalled: '2017-11-07 09:12:08',
                    priority: '7',
                },

            ],
            paging: {
                total: 100,
                currentPage: 1,
                numPerPage: 10,
            },
        };
    }

    reduce (state, action) {
        switch (action.type) {
            case ACTION.APP_INIT:
                return {
                    ...state,
                };

            // case ACTION.REGISTERED_USERS_SUCCESS:
            //     console.log(action.data);
            //     return {
            //         ...state,
            //     };
            //
            // case ACTION.DIALER_START_SUCCESS:
            //     console.log(action.data);
            case ACTION.CHANGE_PRIORITY:
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

            case ACTION.CHANGE_PRIORITY_SUBMIT:
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

            case ACTION.CHANGE_PRIORITY_CANCEL:
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

            default:
                return state;
        };
    };
}

export default new AppStore(appDispatcher);
