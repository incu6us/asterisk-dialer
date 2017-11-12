import React from 'react';
import {Button} from '../Button/Button';
import * as CONSTS from '../../utils/consts';

export default class Table extends React.Component {

    constructor ( props ) {
        super( props );
    }

    render () {
        return (
            <table>
                <thead>
                    {this._makeHead()}
                </thead>
                <tbody>
                    {this._makeBody()}
                </tbody>
            </table>
        )
    }


    _makeHead = () => {
        const {columns} = this.props;
        return <tr>
            {Object.keys( columns ).map( ( column, i ) => (<th key={i}>{columns[ column ]}</th>) )}
        </tr>
    };

    _makeBody = () => {
        const {
            fields,
            columns,
            tableActions,
            isDialer,
            options
        } = this.props;
        const {urls, paging, sortBy, sortOrder} = options;
        const {getCallInProgress, updateCallUrl} = tableActions;
        if (fields.length === 0 && isDialer) {
            return (
                <tr>
                    <td colSpan={Object.keys( columns ).length}>
                        <Button
                            className={'app-button app-button_success'}
                            inscription={'Update table data'}
                            onClick={
                                () => getCallInProgress(
                                    updateCallUrl(
                                        urls[CONSTS.CALL_IN_PROGRESS],
                                        paging.currentPage,
                                        paging.numPerPage,
                                        sortBy,
                                        sortOrder
                                    )
                                )
                            }
                        />
                    </td>
                </tr>
            );
        }
        return fields.map( ( field, i ) =>
            <tr key={i}>
                {Object.keys( columns ).map( ( column, i ) => (
                    <td key={i}>
                        {field[ column ]}
                        {this._makeButtons(field, column)}
                    </td>
                ) )}
            </tr>
        )
    };

    _makeButtons = (field, column) => {
        const {tableActions, options} = this.props;
        const {urls} = options;
        const {deleteRecord} = tableActions;
        switch (column) {
            case CONSTS.DELETE:
                return <Button
                    className={'app-button app-button_delete app-button_delete__small'}
                    inscription={'Delete record'}
                    onClick={() => deleteRecord(urls[CONSTS.CALL_IN_PROGRESS], field.id)}
                />;
            case CONSTS.CHANGE_PRIORITY:
                return <div>
                    {
                        <div className={'app-table-buttons_wrapper'}>
                            <button
                                className={'app-button app-button_success app-button_success__small app-button_up'}
                                onClick={()=>this._handlePriorityChangeDown(field.id, field.priority)}
                                disabled={field.priority === 1}
                            >
                                &#x2b06;
                            </button>
                            <button
                                className={'app-button app-button_delete app-button_success__small'}
                                onClick={()=>this._handlePriorityChangeUp(field.id, field.priority)}
                                disabled={field.priority === 10}
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
        const {tableActions, options} = this.props;
        const {submitPriority} = tableActions;
        const {urls} = options;
        const updatedPriority = priority + 1;
        submitPriority(urls[CONSTS.CALL_IN_PROGRESS], id, updatedPriority);
    };

    _handlePriorityChangeDown = (id, priority) => {
        const {tableActions, options} = this.props;
        const {submitPriority} = tableActions;
        const {urls} = options;
        const updatedPriority = priority - 1;
        submitPriority(urls[CONSTS.CALL_IN_PROGRESS], id, updatedPriority);
    };

}


