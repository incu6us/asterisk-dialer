import React from 'react';
import {Button} from '../Button/Button';
import {DELETE, CHANGE_PRIORITY} from '../../utils/consts';
import {numbers} from "../../utils/validator";

export default class Table extends React.Component {

    constructor ( props ) {
        super( props );

        this.state={
            priority: '',
            isError: false,
        }
    }

    render () {
        const {fields, columns, tableActions} = this.props;
        return (
            <table>
                <thead>
                {this._makeHead( columns )}
                </thead>
                <tbody>
                {this._makeBody( fields, columns, tableActions )}
                </tbody>
            </table>
        )
    }


    _makeHead = (columns) => (
        <tr>
            {Object.keys(columns).map((column, i) => (<th key={i}>{columns[column]}</th>))}
        </tr>
    );

    _makeBody = (fields, columns, tableActions) => (
        fields.map((field, i) =>
            <tr key={i}>
                {Object.keys(columns).map((column, i) => (
                    <td key={i}>
                        {field[column]}
                        {this._makeButtons(field, column, tableActions)}
                    </td>
                ))}
            </tr>
        )
    );

    _makeButtons = (field, column, tableActions) => {
        const {priority, isError} = this.state;
        switch (column) {
            case DELETE:
                return <Button
                    className={'app-button app-button_delete app-button_delete__small'}
                    inscription={'Delete record'}
                    onClick={() => tableActions.deleteRecord(field.id)}
                />;
            case CHANGE_PRIORITY:
                return <div>
                    {
                        <div>
                            <button
                                className={'app-button app-button_success app-button_success__small app-button_up'}
                                onClick={()=>this._handlePriorityChangeUp(field.id, field.priority)}
                                disabled={field.priority === 10}
                            >
                                &#x2b06;
                            </button>
                            <button
                                className={'app-button app-button_delete app-button_success__small'}
                                onClick={()=>this._handlePriorityChangeDown(field.id, field.priority)}
                                disabled={field.priority === 1}
                            >
                                &#x2b07;
                            </button>
                        </div>
                    }
                </div>;
            default:
                return null;
        }
    };

    _handlePriorityChangeUp = (id, priority) => {
        const {tableActions} = this.props;
        const updatedPriority = priority + 1;
        tableActions.submitPriority(id, updatedPriority)
    };

    _handlePriorityChangeDown = (id, priority) => {
        const {tableActions} = this.props;
        const updatedPriority = priority - 1;
        tableActions.submitPriority(id, updatedPriority)
    };

}


