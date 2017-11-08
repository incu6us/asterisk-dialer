import React from 'react';
import {Table} from '../Table/Table';
import {Paging} from '../Paging/Paging';


export const DialerTable = ({orders, columns, pagingChange, paging}) =
>
(
< div >
< Table
fields = {orders}
columns = {columns}
/>
< Paging
baseClassName = {'paging'}
total = {paging.total
}
current = {paging.currentPage
}
numPerPage = {paging.numPerPage
}
onChange = {(page) =
>
pagingChange(page)
}
/>
< /div>
)
;
