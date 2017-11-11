import React from 'react';
import Table from '../Table/Table';
import {Paging} from '../Paging/Paging';
import * as utils from '../../utils/utils';
import * as CONSTS from '../../utils/consts';


export const DialerTable = ({
    orders,
    columns,
    options,
    actions,
    urls
}) => {
    const {paging} = options;
    return <div>
        <Table
            fields={orders}
            columns={columns}
            tableActions={{
                actions,
                updateCallUrl
            }}
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
                ( page ) => actions.pagingChange(
                    updateCallUrl(
                        urls[CONSTS.CALL_IN_PROGRESS],
                        page,
                        20
                    ),
                    page
                )
            }
        />
    </div>
};

const  updateCallUrl = (url, page, limit) => {
    return utils.getUrl(
        url,
        page,
        // sortBy,
        // sortOrder,
        limit
    );
};
