import React from 'react';
import {Button} from '../Button/Button';
import Table from '../Table/Table';
import {DialerTable} from '../DialerTable/DialerTable';
import DropdownList from "../Dropdown/DropdownList";
import * as CONSTS from "../../utils/consts";
import * as utils from '../../utils/utils';


export const Dialer = ({
    operators,
    dialerLists,
    paging,
    pagingChange,
    submitPriority,
    getCallInProgress,
    startDialer,
    stopDialer,
    deleteRecord,
    sortChange,
    limitChange,
    isAppStarted,
    isAppStopped,
    urls,
    sortBy,
    sortOrder,
    clearAll,
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
                    onClick={()=>startDialer(CONSTS.API[CONSTS.START])}
                    isDisabled={isAppStarted}
                />
                <Button
                    className={'app-button app-button_delete right'}
                    inscription={'Stop'}
                    onClick={()=>stopDialer(CONSTS.API[CONSTS.STOP])}
                    isDisabled={isAppStopped}
                />
            </div>
            <Table
                fields={operators}
                columns={CONSTS.COLUMNS}
                tableActions={{}}
                options={{
                    paging,
                    urls,
                    sortBy,
                    sortOrder
                }}
                isDialer={false}
            />
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
            <div className={'app-dropdown-wrapper'}>
                <Button
                    className={'app-button app-button_delete left'}
                    inscription={'Clear all data'}
                    onClick={()=>clearAll(CONSTS.API[CONSTS.CLEAR_ALL])}                    
                />
                <DropdownList
                    items={CONSTS.DROPDOWN_LIST}
                    onChange={
                        (limit)=>limitChange(
                            updateCallUrl(
                                urls[CONSTS.CALL_IN_PROGRESS],
                                paging.currentPage,
                                limit,
                                sortBy,
                                sortOrder,
                            ),
                            limit
                        )
                    }
                    trigger={<button>{ paging.numPerPage } per page</button>}
                    triggerClassName={' button dropdown app-button app-button_secondary app-button_secondary__small'}
                />
            </div>
            <DialerTable
                orders={dialerLists}
                columns={CONSTS.DIALER_COLUMNS}
                options={{
                    paging,
                    sortBy,
                    sortOrder
                }}
                actions={{
                    pagingChange,
                    submitPriority,
                    deleteRecord,
                    getCallInProgress,
                    updateCallUrl,
                    sortChange
                }}
                urls={urls}
            />
        </div>
    );
};

const  updateCallUrl = (url, page, limit, sortBy, sortOrder) => {
    return utils.getUrl(
        url,
        page,
        limit,
        sortBy,
        sortOrder,
    );
};
