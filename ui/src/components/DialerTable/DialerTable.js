import React from 'react';
import Table from '../Table/Table';
import {Paging} from '../Paging/Paging';
import * as CONSTS from '../../utils/consts';


export const DialerTable = ({
    orders,
    columns,
    options,
    actions,
    urls
}) => {
    const {paging} = options;
    const {pagingChange, updateCallUrl} = actions;
    return <div>
        <Table
            fields={orders}
            columns={columns}
            tableActions={actions}
            options={{
                paging,
                urls
            }}
            isDialer={true}
        />
        <Paging
            baseClassName={'paging'}
            total={paging.total}
            current={paging.currentPage}
            numPerPage={paging.numPerPage}
            onChange={
                ( page ) => pagingChange(
                    updateCallUrl(
                        urls[CONSTS.CALL_IN_PROGRESS],
                        page,
                        paging.numPerPage
                    ),
                    page
                )
            }
        />
    </div>
};
