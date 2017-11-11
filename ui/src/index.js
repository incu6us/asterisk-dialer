import React from 'react';
import {render} from 'react-dom';
import App from './components/App';
import css from './styles/style.scss';
import * as actions from './actions/actions';
import * as utils from './utils/utils';
import * as CONSTS from './utils/consts';

const getData = () => {
    actions.getRegisteredUsers(CONSTS.API[CONSTS.REGISTERED_USERS]);
    actions.getCallInProgress(utils.getUrl(CONSTS.API[CONSTS.CALL_IN_PROGRESS], 1, CONSTS.DEFAULT_RECORDS));
    actions.getDialerStatus(CONSTS.API[CONSTS.APP_STATUS]);
};

function init() {
    getData();
    render(<App />,document.getElementById('app'));
}

init();
