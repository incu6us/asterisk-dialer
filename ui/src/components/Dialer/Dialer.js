import React from 'react';
import {Button} from '../Button/Button';
import Table from '../Table/Table';
import {DialerTable} from '../DialerTable/DialerTable';
import * as CONSTS from "../../utils/consts";


export const Dialer = ({
    operators,
    dialerLists,
    paging,
    pagingChange,
    changePriority,
    submitPriority,
    cancelChangePriority,
    getCallInProgress,
    startDialer,
    stopDialer,
    deleteRecord,
    isAppStarted,
    isAppStopped,
}) => {
    return (
        <div className={'app-wrapper'}>
            {
                isAppStarted ?
                    <div className="app-alert-box app-alert-box_success">
                        Application started successfully
                    </div>: null
            }
            {
                isAppStopped ?
                    <div className="app-alert-box app-alert-box_warning">
                        Application stopped successfully
                    </div>: null
            }
            <div className={'app-button-wrapper clearfix'}>
                <Button
                    className={'app-button app-button_success left'}
                    inscription={'Start'}
                    onClick={() => startDialer(CONSTS.API[CONSTS.START])}
                    isDisabled={isAppStarted}
                />
                <Button
                    className={'app-button app-button_delete right'}
                    inscription={'Stop'}
                    onClick={()=>stopDialer(CONSTS.API[CONSTS.STOP])}
                    isDisabled={isAppStopped}
                />
            </div>
            <Table fields={operators} columns={CONSTS.COLUMNS}/>
            {
                paging.total === 0 ?
                    <div className="app-alert-box app-alert-box_warning">
                        You delete all data from Database.
                        Please add data to database and press button ‘Update table data’
                    </div>: null
            }
            {
                dialerLists.length === 0 ?
                    <div className="app-alert-box app-alert-box_warning">
                        You delete some data from app.
                        Please press button ‘Update table data’ and continue working.
                    </div>: null
            }
            <h2>Phone call orders</h2>
            <DialerTable
                orders={dialerLists}
                columns={CONSTS.DIALER_COLUMNS}
                paging={paging}
                actions={{
                    pagingChange,
                    changePriority,
                    submitPriority,
                    cancelChangePriority,
                    deleteRecord,
                    getCallInProgress,
                }}
            />
        </div>
    );
};
