import React from 'react';
import {Button} from '../Button/Button';
import {Table} from '../Table/Table';
import {DialerTable} from '../DialerTable/DialerTable';
import {COLUMNS, DIALER_COLUMNS} from "../../utils/consts";


export const Dialer = ({
    operators,
    dialerLists,
    paging,
    pagingChange,
    changePriority,
    submitPriority,
    cancelChangePriority
}) => {
    return (
        <div className={'app-wrapper'}>
            <div className={'app-button-wrapper clearfix'}>
                <Button className={'app-button app-button_success left'} inscription={'Start'}/>
                <Button className={'app-button app-button_delete right'} inscription={'Stop'}/>
            </div>
            <Table fields={operators} columns={COLUMNS}/>
            <h2>Phone call orders</h2>
            <DialerTable
                orders={dialerLists}
                columns={DIALER_COLUMNS}
                paging={paging}
                actions={{
                    pagingChange,
                    changePriority,
                    submitPriority,
                    cancelChangePriority
                }}
            />
        </div>
    );
};
