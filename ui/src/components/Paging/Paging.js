import React from 'react';
import cx from 'classnames';

/**
 * This file contain implementation of paging component that is a line with links 1, 2, 3... that
 * are a pages. Also it contain a function `paginateList()` that can apply paging to an array.
 *
 * @example
 * <Paging
 *    total={ <number of total items in array> }
 *    current={ <number of current page> }
 *    numPerPage={ <number items per page> }
 *    onChange={ <a handler function> }
 * />
 * <PagingShowing
 *     total={ <number of total items in array> }
 *     current={ <number of current page> }
 *     numPerPage={ <number items per page> }
 * />
 */

export const NUM_PER_PAGE_DEFAULT = 100;

export class Paging extends React.Component {
    /**
     * @returns {Object}
     */
    static get defaultProps() {
        return {
            ...super.defaultProps,
            baseClassName
    :
        'paging',
            startShowed
    :
        3,
            endShowed
    :
        3,
            nearShowed
    :
        1,
            total
    :
        0,
            current
    :
        1,
            numPerPage
    :
        NUM_PER_PAGE_DEFAULT
    }
        ;
    }

    /**
     * @param {Number} page
     * @param {Event} event
     */
    _handlePageButtonClick(page, event) {
        event.preventDefault();
        this.props.onChange(page);
    }

    /**
     * @returns {ReactElement}
     */
    render() {
        const {
            startShowed,
            endShowed,
            nearShowed,
            total,
            current,
            numPerPage,
            baseClassName
        } = this.props;
        const itemsList = [];
        const totalPages = Math.ceil(total / numPerPage);
        if (total <= numPerPage) {
            return null;
        }
        if (current > 1) {
            itemsList.push(
            < li
            key = "-1"
            className = "arrow" >
                < a
            href = "#"
            onClick = {this._handlePageButtonClick.bind(this, current - 1)
        }>&
            laquo;
        <
            /a>
            < /li>
        )
            ;
        }
        for (let i = 0; i < totalPages; i++) {
            if (
                (i + 1 <= startShowed)
                || (i + 1 >= totalPages - endShowed + 1)
                || (current - nearShowed <= i + 1 && i + 1 <= current + nearShowed)
            ) {
                itemsList.push(
                < li
                key = {i}
                className = {cx({'current': current === i + 1
            })
            }>
            <
                a
                href = "#"
                onClick = {this._handlePageButtonClick.bind(this, i + 1)
            }>
                {
                    i + 1
                }
            <
                /a>
                < /li>
            )
                ;
            } else if (
                (i + 1 === current - nearShowed - 1)
                || (i + 1 === current + nearShowed + 1)
                || ((current + nearShowed + 1 <= startShowed || current - nearShowed - 1 >= totalPages - endShowed + 1) && (i + 1 === startShowed + 1))
            ) {
                itemsList.push(
                < li
                key = {i} >
                    < a
                href = "#"
                onClick = {this._handlePageButtonClick.bind(this, i + 1)
            }>&
                hellip;
            <
                /a>
                < /li>
            )
                ;
            }
        }
        if (current < totalPages) {
            itemsList.push(
            < li
            key = "-2"
            className = "arrow" >
                < a
            href = "#"
            onClick = {this._handlePageButtonClick.bind(this, current + 1)
        }>&
            raquo;
        <
            /a>
            < /li>
        )
            ;
        }
        return (
            < ul
        className = {cx('pagination'
    )
    }>
        {
            itemsList
        }
    <
        /ul>
    )
        ;
    }
}

/**
 * Paginate array of objects.
 * It immutable (do not change original array).
 * @param {Array.<Object>} list
 * @param {Number} currentPage
 * @param {Number} numPerPage
 * @return {Array}
 */
export function paginateList(list, currentPage, numPerPage=NUM_PER_PAGE_DEFAULT) {
    const start = (currentPage - 1) * numPerPage;
    const end = (currentPage - 1) * numPerPage + numPerPage;
    return list.filter((item, i) = > i >= start && i < end
)
    ;
}
