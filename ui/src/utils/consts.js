export const REGISTERED_USERS = 'api/v1/dialer/registeredUsers';
export const START = 'api/v1/dialer/start';
export const STOP = 'api/v1/dialer/stop';
export const APP_STATUS = 'api/v1/dialer/status';
export const CALL_IN_PROGRESS = 'api/v1/dialer/msisdn/inProgress';
export const CLEAR_ALL = '/api/v1/dialer/msisdn/clearAll';

const ACTION = 'action';
const EXTEN = 'exten';
const PEER = 'peer';
const PEER_STATUS = 'peerStatus';

const ID = 'id';
const MSISDN = 'msisdn';
const STATUS = 'status';
const TIME = 'time';
const CAUSE_TXT = 'causeTxt';
const EVENT = 'event';
const CALLER_ID_NUM = 'callerIdNum';
const TIME_CALLED = 'timeCalled';
export const PRIORITY = 'priority';
export const DELETE = 'delete';
export const CHANGE_PRIORITY = 'changePriority';

export const COLUMNS = {
    [ACTION]: 'Operator Action',
    [EXTEN]: 'Phone number',
    [PEER]: 'Operator',
    [PEER_STATUS]: 'Operator Status'
};

export const DIALER_COLUMNS = {
    [ID]: 'ID',
    [MSISDN]: 'Phone Number',
    [STATUS]: 'Status',
    [TIME]: 'Time',
    [CAUSE_TXT]: 'Cause Txt',
    [EVENT]: 'Event',
    [CALLER_ID_NUM]: 'Operator',
    [TIME_CALLED]: 'Time Called',
    [PRIORITY]: 'Priority',
    [CHANGE_PRIORITY]: '',
    [DELETE]: '',
};

export const DROPDOWN_LIST = [
    {
        text: '20 per page',
        value: 20,
    },
    {
        text: '50 per page',
        value: 50,
    },
    {
        text: '100 per page',
        value: 100
    }
];

export const ASC = 'asc';
export const DESC = 'desc';

export const getHostFn = () => {
    return window.location.protocol+"//"+window.location.host+"/{API}";
};

export const DEFAULT_RECORDS = 20;

export const API = {
    [START]: getHostFn().replace('{API}', START),
    [STOP]: getHostFn().replace('{API}', STOP),
    [APP_STATUS]: getHostFn().replace('{API}', APP_STATUS),
    [REGISTERED_USERS]: getHostFn().replace('{API}', REGISTERED_USERS),
    [CALL_IN_PROGRESS]: getHostFn().replace('{API}', CALL_IN_PROGRESS),
    [CLEAR_ALL]: getHostFn().replace('{API}', CLEAR_ALL)
};
