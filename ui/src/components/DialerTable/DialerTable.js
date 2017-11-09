import React from 'react';
import Table from '../Table/Table';
import {Paging} from '../Paging/Paging';


export const DialerTable = ({
    orders,
    columns,
    paging,
    actions,
}) => {
    return <div>
        <Table
            fields={orders}
            columns={columns}
            tableActions={actions}
        />
        <Paging
            baseClassName={'paging'}
            total={paging.total}
            current={paging.currentPage}
            numPerPage={paging.numPerPage}
            onChange={( page ) => actions.pagingChange( page )}
        />
    </div>
};
