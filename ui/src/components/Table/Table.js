import React from 'react';
import {Button} from '../Button/Button';
import {DELETE, CHANGE_PRIORITY} from '../../utils/consts';

export const Table = ({fields, columns}) =
>
{
    return (
        < table >
        < thead >
        {makeHead(columns)}
        < /thead>
        < tbody >
        {makeBody(fields, columns)}
        < /tbody>
        < /table>
)
}
;


const makeHead = (columns) =
>
(
< tr >
{Object.keys(columns).map(column = > ( < th > {columns[column]} < /th>))}
< /tr>
)
;

const makeBody = (fields, columns) =
>
(
    fields.map(field = >
    < tr >
    {Object.keys(columns).map(column = > (
    < td >
    {field[column]}
{
    makeButtons(column)
}
<
/td>
))
}
<
/tr>
)
)
;

const makeButtons = (column) =
>
{
    switch (column) {
        case DELETE:
            return
        <
            Button
            className = {'app-button app-button_delete app-button_delete__small'}
            inscription = {'Delete order'}
            />;
        case CHANGE_PRIORITY:
            return
        <
            div >
            < Button
            className = {'app-button app-button_success app-button_success__small'}
            inscription = {'Change Priority'}
            />
            < label
            htmlFor = "" >
                Add
            new priority
            < input
            type = "text"
            value = ''
            name = 'changePriority' / >
                < /label>
                < /div>;
        default:
            return null;
    }
}
;
