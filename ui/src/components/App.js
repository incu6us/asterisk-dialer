import React, { Component } from 'react';
import { Container } from 'flux/utils';
import {Dialer} from './Dialer/Dialer';
// store
import appStore from '../stores/store';
// actions
import * as actions from '../actions/actions'
import * as CONSTS from '../utils/consts';

class App extends Component {

    static getStores() {
        return [ appStore ];
    }

    static calculateState(prevState, props) {
        return {
            ...appStore.getState()
        }
    }

    componentDidMount () {
        setInterval(()=> actions.getRegisteredUsers(CONSTS.API[CONSTS.REGISTERED_USERS]), 10000);
        setInterval(()=> actions.getDialerStatus(CONSTS.API[CONSTS.APP_STATUS]), 1000);
    }

    render() {
        console.log(this.state);
        return (
            <Dialer
                {...actions}
                {...this.state}
            />
        );
    }
}

export default new Container.create(App);
