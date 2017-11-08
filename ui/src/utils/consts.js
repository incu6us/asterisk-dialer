export const registeredUSers = 'registeredUsers';
export const start = 'start';
export const stop = 'stop';
export const status = 'status';

const ACTION = 'action';
const EXTEN = 'exten';
const PEER = 'peer';
const PEER_STATUS = 'peerStatus';

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

export const getHostFn = () => {
    return window.location.protocol+"//"+window.location.host+"/api/v1/dialer/{API}";
};
