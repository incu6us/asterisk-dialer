import React from 'react';
import {Button} from '../Button/Button';
import {DELETE, CHANGE_PRIORITY} from '../../utils/consts';

export const Table = ({fields, columns, tableActions}) => {
    return (
        <table>
            <thead>
                {makeHead(columns)}
            </thead>
            <tbody>
                {makeBody(fields, columns, tableActions)}
            </tbody>
        </table>
    )
};


const makeHead = (columns) => (
    <tr>
        {Object.keys(columns).map((column, i) => (<th key={i}>{columns[column]}</th>))}
    </tr>
);

const makeBody = (fields, columns, tableActions) => (
    fields.map((field, i) =>
        <tr key={i}>
            {Object.keys(columns).map((column, i) => (
                <td key={i}>
                    {field[column]}
                    {makeButtons(field, column, tableActions)}
                </td>
            ))}
        </tr>
    )
);

const makeButtons = (field, column, tableActions) => {
    switch (column) {
        case DELETE:
            return <Button
                className={'app-button app-button_delete app-button_delete__small'}
                inscription={'Delete order'}
            />;
        case CHANGE_PRIORITY:
            return <div>
                {
                    field.isChanging ?
                        <div>
                            <label htmlFor="">
                                Add new priority
                                <input
                                    type="text"
                                    value=''
                                    name='changePriority'
                                    onChange={() => tableActions.typingPriority(field.msisdn)}
                                />
                            </label>
                            <Button
                                className={'app-button app-button_success app-button_success__small'}
                                inscription={'Submit Priority'}
                                onClick={()=>tableActions.submitPriority(field.msisdn)}
                            />
                            <Button
                                className={'app-button app-button_alert app-button_alert__small'}
                                inscription={'Cancel'}
                                onClick={()=>tableActions.cancelChangePriority(field.msisdn)}
                            />
                        </div>:
                        <Button
                            className={'app-button app-button_success app-button_success__small'}
                            inscription={'Change Priority'}
                            onClick={() => tableActions.changePriority(field.msisdn)}
                        />
                }
            </div>;
        default:
            return null;
    }
};
