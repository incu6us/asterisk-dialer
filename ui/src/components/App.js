import React, { Component } from 'react';
import { Container } from 'flux/utils';
import {Dialer} from './Dialer/Dialer';



// store
import appStore from '../stores/store';
// actions
import * as actions from '../actions/actions'

class App extends Component {

    static getStores() {
        return [ appStore ];
    }

    static calculateState(prevState, props) {
        return {
            ...appStore.getState()
        }
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
