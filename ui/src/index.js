import React from 'react';
import {render} from 'react-dom';
import App from './components/App';
import css from './styles/style.scss';
import * as actions from './actions/actions'
import * as consts from './utils/consts';

const getData = () =
>
{
    const registeredUsersUrl = consts.getHostFn().replace('{API}', consts.registeredUSers);
    const dialerStatusUrl = consts.getHostFn().replace('{API}', consts.status);
    //actions.getRegisteredUsers(registeredUsersUrl);
    //actions.getDialerStatus(dialerStatusUrl);
}
;

function init() {
    getData();
    render( < App / >, document.getElementById('app')
)
    ;
}

init();
