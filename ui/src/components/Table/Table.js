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
                        field.isChanging ?
                            <div>
                                <label htmlFor="">
                                    Add new priority
                                    <input
                                        type="text"
                                        value={priority}
                                        name='changePriority'
                                        className={isError ? 'app-input_hasError': ''}
                                        onChange={(e) => this._handlePriorityChange(e)}
                                    />
                                    {
                                        isError ?
                                            <small className={'help-text_error'}>
                                                Please enter number 1 - 10
                                            </small> : null
                                    }
                                </label>
                                <Button
                                    className={'app-button app-button_success app-button_success__small'}
                                    inscription={'Submit'}
                                    onClick={()=>this._handlePrioritySubmit(field.id, priority)}
                                />
                                <Button
                                    className={'app-button app-button_alert app-button_alert__small'}
                                    inscription={'Cancel'}
                                    onClick={()=>tableActions.cancelChangePriority(field.id)}
                                />
                            </div>:
                            <Button
                                className={'app-button app-button_success app-button_success__small'}
                                inscription={'Change Priority'}
                                onClick={()=>tableActions.changePriority(field.id)}
                            />
                    }
                </div>;
            default:
                return null;
        }
    };

    _handlePriorityChange = ({target}) => {
        const isValid = numbers(target.value);
        if ((isValid && target.value <= 10) || !target.value) {
            this.setState(state => ({
                ...state,
                priority: target.value,
                isError: false,
            }));
        } else {
            this.setState(state => ({
                ...state,
                isError: true
            }));
        }
    };

    _handlePrioritySubmit = (id, priority) => {
        const {tableActions} = this.props;
        tableActions.submitPriority(id, Number(priority))
            .then(() => this.setState(state => ({
                ...state,
                priority: '',
            })))
    }
}


