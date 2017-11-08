import React from 'react';
import {render} from 'react-dom';
import App from './components/App';
import css from './styles/style.scss';
import * as actions from './actions/actions'
import * as CONSTS from './utils/consts';

const getData = () => {
    actions.getRegisteredUsers(CONSTS.API[CONSTS.REGISTERED_USERS]);
    //actions.getDialerStatus(dialerStatusUrl);
};

function init() {
    getData();
    render(<App />,document.getElementById('app'));
}

init();
