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

    shouldComponentUpdate (newProps, newState) {
        console.log(newState.dialerLists.length);
        return newState.dialerLists.length === 0;
    }

    componentDidMount(){
        const {dialerLists} = this.state;
        setInterval(()=> actions.getRegisteredUsers(CONSTS.API[CONSTS.REGISTERED_USERS]), 10000);
        setInterval(()=> actions.getDialerStatus(CONSTS.API[CONSTS.APP_STATUS]), 1000);
        let timer = null;
        console.log('length ->>', dialerLists.length);
        if (dialerLists.length === 0) {
            timer = setInterval(()=>actions.getCallInProgress(CONSTS.API[CONSTS.CALL_IN_PROGRESS]), 2000);
        }
        clearInterval(timer)
    }

    render() {
        return (
            <Dialer
                {...actions}
                {...this.state}
            />
        );
    }
}

export default new Container.create(App);
